package fuzzer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/hideckies/fuzzagotchi/cmd"
	"github.com/hideckies/fuzzagotchi/pkg/output"
	"github.com/hideckies/fuzzagotchi/pkg/util"
	"github.com/schollz/progressbar/v3"

	"github.com/fatih/color"
)

type Config struct {
	Context    context.Context `json:"context"`
	Cookie     string          `json:"cookie"`
	Delay      string          `json:"delay"`
	EGG        bool            `json:"egg"`
	Header     string          `json:"header"`
	Method     string          `json:"method"`
	PostData   string          `json:"post_data"`
	Retry      int             `json:"retry"`
	StatusCode []int           `json:"match_status"`
	Threads    int             `json:"threads"`
	Timeout    int             `json:"timeout"`
	URL        string          `json:"url"`
	UserAgent  string          `json:"user_agent"`
	Verbose    bool            `json:"verbose"`
	Wordlist   string          `json:"wordlist"`
}

type Fuzzer struct {
	Config    Config     `json:"config"`
	Request   Request    `json:"request"`
	Responses []Response `json:"response"`

	TotalWords int `json:"total_words"`

	// mu *sync.Mutex `json:"-"`
}

// Initialize a new Fuzzer
func NewFuzzer(ctx context.Context, options cmd.CmdOptions, totalWords int) Fuzzer {
	var f Fuzzer

	egg := false
	if strings.Contains(options.URL, "EGG") {
		egg = true
	}

	f.Config = Config{
		Cookie:     options.Cookie,
		Context:    ctx,
		Delay:      options.Delay,
		EGG:        egg,
		Header:     options.Header,
		Method:     options.Method,
		PostData:   options.PostData,
		Retry:      options.Retry,
		StatusCode: options.StatusCode,
		Threads:    options.Threads,
		Timeout:    options.Timeout,
		URL:        options.URL,
		UserAgent:  options.UserAgent,
		Verbose:    options.Verbose,
		Wordlist:   options.Wordlist,
	}
	f.Request = Request{}
	f.TotalWords = totalWords
	return f
}

// Run to fuzz
func (f *Fuzzer) Run() {
	runtime.GOMAXPROCS(f.Config.Threads)
	var wg sync.WaitGroup

	bar := *output.NewProgressBar(f.TotalWords, "Fuzzing...")

	readFile, err := os.Open(f.Config.Wordlist)
	if err != nil {
		log.Fatal(err)
	}
	defer readFile.Close()

	scanner := bufio.NewScanner(readFile)
	scanner.Split(bufio.ScanLines)

	wordCh := make(chan string, f.Config.Threads)
	respCh := make(chan Response, f.Config.Threads)

	for i := 0; i < f.Config.Threads; i++ {
		wg.Add(1)
		go f.worker(&wg, i, wordCh, respCh, bar)
	}

	for scanner.Scan() {
		bar.Add(1)
		wordCh <- scanner.Text()
	}

	for len(respCh) > 0 {
		r := <-respCh
		f.addResponse(r)
	}

	bar.Close()

	// Output result
	f.output()

	close(wordCh)
	wg.Wait()
	close(respCh)
}

// Worker to fuzz using given word
func (f *Fuzzer) worker(wg *sync.WaitGroup, id int, wordCh chan string, respCh chan Response, bar progressbar.ProgressBar) {
	defer wg.Done()

	for word := range wordCh {
		resp, err := f.process(word, 1)
		if err != nil {
			if f.Config.Verbose {
				fmt.Println(err)
			}
		}
		// respCh <- resp

		f.addResponse(resp)
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
		// if f.Config.Retry > n {
		// 	time.Sleep(getDelay(f.Config.Delay))
		// 	resp, err := f.process(word, n+1)
		// }
		// return Response{}, fmt.Errorf("%s", err)
	}

	return resp, nil
}

// Adjust response
func (f *Fuzzer) addResponse(resp Response) {
	if util.ContainInt(f.Config.StatusCode, resp.StatusCode) && resp.ContentLength >= 0 {
		f.Responses = append(f.Responses, resp)
	}
}

// Print results
func (f *Fuzzer) output() {
	tw := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', tabwriter.TabIndent)
	defer tw.Flush()

	fmt.Fprintf(tw, "%s\n", color.YellowString(output.TMPL_BAR_SINGLE_M))
	fmt.Fprintf(tw,
		"%s %s\t%s\t%s\n",
		color.CyanString("+"),
		color.CyanString("WORD"),
		color.YellowString("Status Code"),
		color.HiMagentaString("Content Length"))
	fmt.Fprintf(tw, "%s\n", color.YellowString(output.TMPL_BAR_SINGLE_M))

	// result := color.CyanString(
	// 	"%-10s\t\tStatus: %d, Content Length: %d, Duration: %.2fs",
	// 	resp.Word,
	// 	resp.StatusCode,
	// 	resp.ContentLength,
	// 	resp.Delay.Abs().Seconds())

	for _, resp := range f.Responses {
		fmt.Fprintf(tw,
			"%s\t%s\t%s\n",
			color.CyanString(resp.Word),
			color.YellowString("%d", resp.StatusCode),
			color.HiMagentaString("%d", resp.ContentLength))
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
