package fuzzer

import (
	"bufio"
	"context"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/cmd"
)

type Config struct {
	Context   context.Context `json:"context"`
	Cookie    string          `json:"cookie"`
	Delay     string          `json:"delay"`
	EGG       bool            `json:"egg"`
	Header    string          `json:"header"`
	Method    string          `json:"method"`
	PostData  string          `json:"post_data"`
	Retry     int             `json:"retry"`
	Threads   int             `json:"threads"`
	Timeout   int             `json:"timeout"`
	URL       string          `json:"url"`
	UserAgent string          `json:"user_agent"`
	Verbose   bool            `json:"verbose"`
	Wordlist  string          `json:"wordlist"`
}

type Fuzzer struct {
	Config  Config  `json:"config"`
	Request Request `json:"request"`
}

// Initialize a new Fuzzer
func NewFuzzer(options cmd.CmdOptions, ctx context.Context) Fuzzer {
	var f Fuzzer

	egg := false
	if strings.Contains(options.URL, "EGG") {
		egg = true
	}

	f.Config = Config{
		Cookie:    options.Cookie,
		Context:   ctx,
		Delay:     options.Delay,
		EGG:       egg,
		Header:    options.Header,
		Method:    options.Method,
		PostData:  options.PostData,
		Retry:     options.Retry,
		Threads:   options.Threads,
		Timeout:   options.Timeout,
		URL:       options.URL,
		UserAgent: options.UserAgent,
		Verbose:   options.Verbose,
		Wordlist:  options.Wordlist,
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

	for i := 0; i < f.Config.Threads; i++ {
		go f.worker(&wg, i, wordCh)
	}

	for fileScanner.Scan() {
		word := fileScanner.Text()
		wordCh <- word
	}
	close(wordCh)
	readFile.Close()
	wg.Wait()
}

// Worker to fuzz using given word
func (f *Fuzzer) worker(wg *sync.WaitGroup, id int, wordCh chan string) {
	defer wg.Done()
	for word := range wordCh {
		f.process(word, 1)
		time.Sleep(getDelay(f.Config.Delay))
	}
}

// Process to send a request
func (f *Fuzzer) process(word string, n int) {
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
		return
	}
	f.output(resp)
}

// Print results
func (f *Fuzzer) output(resp Response) {
	result := color.CyanString(
		"%-10s\t\tStatus: %d, Content Length: %d, Duration: %.2fs",
		resp.Word,
		resp.StatusCode,
		resp.ContentLength,
		resp.Delay.Abs().Seconds())
	resultFailed := color.RedString("[x] %v", result)

	rcl, _ := regexp.Compile("^([1-9][0-9]*|0)$")
	rclrange, _ := regexp.Compile("^(([1-9][0-9]*|0)-([1-9][0-9]*|0))$")

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
