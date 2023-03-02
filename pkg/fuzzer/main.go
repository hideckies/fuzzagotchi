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
	"time"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/cmd"
	"github.com/hideckies/fuzzagotchi/pkg/output"
	"github.com/hideckies/fuzzagotchi/pkg/util"
	"github.com/schollz/progressbar/v3"
)

type Config struct {
	Context        context.Context `json:"context"`
	Cookie         string          `json:"cookie"`
	Delay          string          `json:"delay"`
	FollowRedirect bool            `json:"follow_redirect"`
	FuzzType       string          `json:"fuzztype"`
	EGG            bool            `json:"egg"`
	Header         string          `json:"header"`
	HideLength     string          `json:"hide_length"`
	HideStatus     []int           `json:"hide_status"`
	Host           string          `json:"host"`
	MatchLength    string          `json:"match_length"`
	MatchStatus    []int           `json:"match_status"`
	Method         string          `json:"method"`
	PostData       string          `json:"post_data"`
	Retry          int             `json:"retry"`
	Threads        int             `json:"threads"`
	Timeout        int             `json:"timeout"`
	URL            string          `json:"url"`
	UserAgent      string          `json:"user_agent"`
	Verbose        bool            `json:"verbose"`
	Wordlist       string          `json:"wordlist"`
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
func NewFuzzer(ctx context.Context, options cmd.CmdOptions, fuzztype string, totalWords int) Fuzzer {
	var f Fuzzer

	egg := false
	if strings.Contains(options.URL, "EGG") {
		egg = true
	}

	f.Config = Config{
		Cookie:         options.Cookie,
		Context:        ctx,
		Delay:          options.Delay,
		EGG:            egg,
		FuzzType:       fuzztype,
		FollowRedirect: options.FollowRedirect,
		Header:         options.Header,
		HideLength:     options.HideLength,
		HideStatus:     options.HideStatus,
		Host:           extractHost(options.URL),
		MatchLength:    options.MatchLength,
		MatchStatus:    options.MatchStatus,
		Method:         options.Method,
		PostData:       options.PostData,
		Retry:          options.Retry,
		Threads:        options.Threads,
		Timeout:        options.Timeout,
		URL:            options.URL,
		UserAgent:      options.UserAgent,
		Verbose:        options.Verbose,
		Wordlist:       options.Wordlist,
	}

	// Auto EGG
	if fuzztype == "" {
		f.Config.URL = util.AdjustUrlSuffix(f.Config.URL) + "EGG"
	}

	f.Request = NewRequest(f.Config)
	f.Responses = make([]Response, 0)
	// f.ResponsePool = make([]Response, 0)
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
		go f.worker(&wg, i, wordCh, bar)
	}

	for scanner.Scan() {
		bar.Add(1)

		word := scanner.Text()
		wordCh <- word

		// Auto EGG
		if f.Config.FuzzType == "" {
			// TXT files
			wTxt := word + ".txt"
			wordCh <- wTxt
			// HTML files
			wHtml := word + ".html"
			wordCh <- wHtml
			// PHP files
			wPhp := word + ".php"
			wordCh <- wPhp
			// Hidden files
			wHidden := "." + word
			wordCh <- wHidden
		}

	}

	bar.Close()
	close(wordCh)
	wg.Wait()

	fmt.Println()

	f.printResultHeader()

	// Finding information in each page
	// explore := NewExplore(f.Responses)
	// explore.explore()
}

// Worker to fuzz using given word
func (f *Fuzzer) worker(wg *sync.WaitGroup, id int, wordCh chan string, bar progressbar.ProgressBar) {
	defer wg.Done()

	for word := range wordCh {
		resp, err := f.process(word, 1)
		if err != nil {
			if f.Config.Verbose {
				fmt.Println(err)
			}
		}

		if f.matcher(resp) {
			f.Responses = append(f.Responses, resp)
			f.printResultURL(resp, bar)
		}
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

// Check if the response matches rules
func (f *Fuzzer) matcher(resp Response) bool {
	if util.ContainInt(f.Config.MatchStatus, resp.StatusCode) && !util.ContainInt(f.Config.HideStatus, resp.StatusCode) && f.matchContentLength(resp.ContentLength) {
		return true
	}
	return false
}

// Print result of URL fuzzing
func (f *Fuzzer) printResultURL(resp Response, bar progressbar.ProgressBar) {
	if f.Config.FuzzType != "" && f.Config.FuzzType != "URL" {
		return
	}

	bar.Clear()

	result := fmt.Sprintf("%-32s %s%s %s%s",
		color.CyanString(resp.Path),
		color.YellowString("("),
		color.GreenString("status: %d", resp.StatusCode),
		color.HiMagentaString("length: %d", resp.ContentLength),
		color.YellowString(")"))

	if resp.RedirectPath != "" {
		fmt.Printf("%s -> %s\n", result, color.GreenString(resp.RedirectPath))
	} else {
		fmt.Printf("%s\n", result)
	}
}

// Print result of header fuzzing
func (f *Fuzzer) printResultHeader() {
	if f.Config.FuzzType != "Header" {
		return
	}

	color.Yellow(output.TMPL_BAR_DOUBLE_M)
	fmt.Printf("%s %s\n", color.CyanString("+"), color.CyanString("HEADER FUZZING"))
	color.Yellow(output.TMPL_BAR_DOUBLE_M)

	if len(f.Responses) > 0 {
		lengthToWords := make(map[int][]string)
		for _, resp := range f.Responses {
			lengthToWords[resp.ContentLength] = append(lengthToWords[resp.ContentLength], resp.Word)
		}

		// Exclude
		maxCnt := 0
		keyOfMaxCnt := 0
		for key, val := range lengthToWords {
			cnt := len(val)
			if maxCnt < cnt {
				maxCnt = cnt
				keyOfMaxCnt = key
			}
		}
		delete(lengthToWords, keyOfMaxCnt)

		if len(lengthToWords) > 0 {
			// Print result
			for _, val := range lengthToWords {
				for _, v := range val {
					fmt.Printf("%-23s",
						color.CyanString(v))
				}
			}
		} else {
			fmt.Println("No result found")
		}
	} else {
		fmt.Println("No result found.")
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
