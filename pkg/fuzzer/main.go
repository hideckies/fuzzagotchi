package fuzzer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/hideckies/fuzzagotchi/cmd"
	"github.com/hideckies/fuzzagotchi/pkg/util"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

type Config struct {
	Context     context.Context `json:"context"`
	Cookie      string          `json:"cookie"`
	Delay       string          `json:"delay"`
	EGG         bool            `json:"egg"`
	Header      string          `json:"header"`
	MatchStatus []int           `json:"match_status"`
	Method      string          `json:"method"`
	PostData    string          `json:"post_data"`
	Retry       int             `json:"retry"`
	Threads     int             `json:"threads"`
	Timeout     int             `json:"timeout"`
	URL         string          `json:"url"`
	UserAgent   string          `json:"user_agent"`
	Verbose     bool            `json:"verbose"`
	Wordlist    string          `json:"wordlist"`
}

type Fuzzer struct {
	Config  Config  `json:"config"`
	Request Request `json:"request"`

	TotalWords  int                     `json:"total_words"`
	ProgressBar progressbar.ProgressBar `json:"progressbar"`
}

// Initialize a new Fuzzer
func NewFuzzer(options cmd.CmdOptions, ctx context.Context) Fuzzer {
	var f Fuzzer

	egg := false
	if strings.Contains(options.URL, "EGG") {
		egg = true
	}

	f.Config = Config{
		Cookie:      options.Cookie,
		Context:     ctx,
		Delay:       options.Delay,
		EGG:         egg,
		Header:      options.Header,
		MatchStatus: options.MatchStatus,
		Method:      options.Method,
		PostData:    options.PostData,
		Retry:       options.Retry,
		Threads:     options.Threads,
		Timeout:     options.Timeout,
		URL:         options.URL,
		UserAgent:   options.UserAgent,
		Verbose:     options.Verbose,
		Wordlist:    options.Wordlist,
	}
	f.Request = Request{}
	return f
}

// Run to fuzz
func (f *Fuzzer) Run() {
	var wg sync.WaitGroup
	wg.Add(f.Config.Threads)

	readFile, err := os.Open(f.Config.Wordlist)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	wordCh := make(chan string, f.Config.Threads)
	respCh := make(chan Response, f.Config.Threads)

	for i := 0; i < f.Config.Threads; i++ {
		go f.worker(&wg, i, wordCh, respCh)
	}

	for fileScanner.Scan() {
		f.TotalWords++
		word := fileScanner.Text()
		wordCh <- word
	}

	// progBar := *output.NewProgressBar(f.TotalWords, "Fuzzing...")

	for len(respCh) > 0 {
		r := <-respCh
		f.output(r)
	}

	close(wordCh)
	readFile.Close()
	wg.Wait()
	close(respCh)
}

// Worker to fuzz using given word
func (f *Fuzzer) worker(wg *sync.WaitGroup, id int, wordCh chan string, respCh chan Response) {
	defer wg.Done()

	for word := range wordCh {
		resp, err := f.process(word, 1)
		if err == nil {
			respCh <- resp
			// if the respCh is max...
			if len(respCh) == f.Config.Threads {
				for j := 0; j < f.Config.Threads; j++ {
					r := <-respCh
					f.output(r)
				}
			}
		}
		time.Sleep(getDelay(f.Config.Delay))
	}

}

// Process to send a request
func (f *Fuzzer) process(word string, n int) (Response, error) {
	req := NewRequest(f.Config)
	// Send request
	resp, err := req.Send(word)
	if err != nil {
		if f.Config.Verbose {
			color.Red("%-10s\t\tError: %s\n", word, err)
		}
		// Retry to send a request until reaching the retry limit.
		if f.Config.Retry > n {
			time.Sleep(getDelay(f.Config.Delay))
			f.process(word, n+1)
		}
		return Response{}, fmt.Errorf("%s", err)
	}
	return resp, nil
}

// Print results
func (f *Fuzzer) output(resp Response) {
	// result := color.CyanString(
	// 	"%-10s\t\tStatus: %d, Content Length: %d, Duration: %.2fs",
	// 	resp.Word,
	// 	resp.StatusCode,
	// 	resp.ContentLength,
	// 	resp.Delay.Abs().Seconds())

	if util.ContainInt(f.Config.MatchStatus, resp.StatusCode) && resp.ContentLength > 0 {
		// fmt.Fprintf(&f.TabWriter, "%s\t%s\n", color.CyanString(resp.Word), color.YellowString("%d", resp.StatusCode))
		fmt.Printf("%s\t%s\n", color.CyanString(resp.Word), color.YellowString("%d", resp.StatusCode))
	}

	// rcl, _ := regexp.Compile("^([1-9][0-9]*|0)$")
	// rclrange, _ := regexp.Compile("^(([1-9][0-9]*|0)-([1-9][0-9]*|0))$")

	// if rcl.MatchString(resp.Config.NoContentLength) {
	// 	ncl, _ := strconv.Atoi(resp.Config.NoContentLength)
	// 	if util.ContainInt(resp.Config.Status, resp.StatusCode) && ncl != resp.ContentLength {
	// 		fmt.Println(result)
	// 	} else if resp.Config.Verbose {
	// 		fmt.Println(resultFailed)
	// 	}
	// } else {
	// 	if rclrange.MatchString(resp.Config.ContentLength) {
	// 		contentlengths := strings.Split(resp.ContentLength, "-")
	// 		cmin, _ := strconv.Atoi(contentlengths[0])
	// 		cmax, _ := strconv.Atoi(contentlengths[1])
	// 		if util.ContainInt(resp.Status, resp.StatusCode) && (cmin <= resp.ContentLength && resp.ContentLength <= cmax) {
	// 			fmt.Println(result)
	// 		} else if resp.Config.Verbose {
	// 			fmt.Println(resultFailed)
	// 		}
	// 	} else if rcl.MatchString(resp.Config.ContentLength) {
	// 		cl, _ := strconv.Atoi(resp.Config.ContentLength)
	// 		if util.ContainInt(resp.Config.Status, resp.StatusCode) && cl == resp.ContentLength {
	// 			fmt.Println(result)
	// 		} else if resp.Config.Verbose {
	// 			fmt.Println(resultFailed)
	// 		}
	// 	} else {
	// 		if util.ContainInt(resp.Config.Status, resp.StatusCode) {
	// 			fmt.Println(result)
	// 		} else if resp.Config.Verbose {
	// 			fmt.Println(resultFailed)
	// 		}
	// 	}
	// }
}
