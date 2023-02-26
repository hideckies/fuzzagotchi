package fuzzer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"strings"
	"sync"
	"text/tabwriter"
	"time"

	"github.com/hideckies/fuzzagotchi/cmd"
	"github.com/hideckies/fuzzagotchi/pkg/output"
	"github.com/hideckies/fuzzagotchi/pkg/util"

	"github.com/fatih/color"
)

type Config struct {
	ContentLength string          `json:"content_length"`
	Context       context.Context `json:"context"`
	Cookie        string          `json:"cookie"`
	Delay         string          `json:"delay"`
	EGG           bool            `json:"egg"`
	Header        string          `json:"header"`
	Host          string          `json:"host"`
	Method        string          `json:"method"`
	PostData      string          `json:"post_data"`
	Retry         int             `json:"retry"`
	StatusCode    []int           `json:"match_status"`
	Threads       int             `json:"threads"`
	Timeout       int             `json:"timeout"`
	URL           string          `json:"url"`
	UserAgent     string          `json:"user_agent"`
	Verbose       bool            `json:"verbose"`
	Wordlist      string          `json:"wordlist"`
}

type Fuzzer struct {
	Config    Config     `json:"config"`
	Request   Request    `json:"request"`
	Responses []Response `json:"response"`

	TotalWords int      `json:"total_words"`
	ErrorWords []string `json:"error_words"`

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
		ContentLength: options.ContentLength,
		Cookie:        options.Cookie,
		Context:       ctx,
		Delay:         options.Delay,
		EGG:           egg,
		Header:        options.Header,
		Host:          extractHost(options.URL),
		Method:        options.Method,
		PostData:      options.PostData,
		Retry:         options.Retry,
		StatusCode:    options.StatusCode,
		Threads:       options.Threads,
		Timeout:       options.Timeout,
		URL:           options.URL,
		UserAgent:     options.UserAgent,
		Verbose:       options.Verbose,
		Wordlist:      options.Wordlist,
	}
	f.Request = NewRequest(f.Config)
	f.TotalWords = totalWords
	f.ErrorWords = make([]string, 0)
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

	for i := 0; i < f.Config.Threads; i++ {
		wg.Add(1)
		go f.worker(&wg, i, wordCh)
	}

	for scanner.Scan() {
		bar.Add(1)
		wordCh <- scanner.Text()
	}

	bar.Close()
	close(wordCh)
	wg.Wait()

	// Output result
	f.output()
}

// Worker to fuzz using given word
func (f *Fuzzer) worker(wg *sync.WaitGroup, id int, wordCh chan string) {
	defer wg.Done()

	for word := range wordCh {
		resp, err := f.process(word, 1)
		if err != nil {
			if f.Config.Verbose {
				fmt.Println(err)
			}
		}
		f.addResponse(resp)
		time.Sleep(getDelay(f.Config.Delay))
	}
}

// Process to send a request
func (f *Fuzzer) process(word string, n int) (Response, error) {
	var resp Response
	var err error

	cnt := 0
	ok := false
	for !ok && cnt <= f.Config.Retry {
		resp, err = f.Request.Send(word)
		if err == nil {
			ok = true
		}
		time.Sleep(getDelay(f.Config.Delay))
		cnt++
	}
	return resp, nil
}

// Adjust response
func (f *Fuzzer) addResponse(resp Response) {
	if util.ContainInt(f.Config.StatusCode, resp.StatusCode) && resp.ContentLength >= 0 && f.matchContentLength(resp.ContentLength) {
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

	for _, resp := range f.Responses {
		fmt.Fprintf(tw,
			"%s\t%s\t%s\n",
			color.CyanString(resp.Word),
			color.YellowString("%d", resp.StatusCode),
			color.HiMagentaString("%d", resp.ContentLength))
	}
}

// Extract hostname from URL
func extractHost(_url string) string {
	u, err := url.Parse(_url)

	if err != nil {
		return ""
	}
	return u.Hostname()
}
