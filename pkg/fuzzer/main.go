package fuzzer

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"runtime"
	"strconv"
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
	HideWords      string          `jsong:"hide_words"`
	Host           string          `json:"host"`
	MatchLength    string          `json:"match_length"`
	MatchStatus    []int           `json:"match_status"`
	MatchWords     string          `json:"match_words"`
	Method         string          `json:"method"`
	PostData       string          `json:"post_data"`
	Proxy          string          `jsong:"proxy"`
	Retry          int             `json:"retry"`
	Scan           bool            `json:"scan"`
	Threads        int             `json:"threads"`
	Timeout        int             `json:"timeout"`
	URL            string          `json:"url"`
	UserAgent      string          `json:"user_agent"`
	Verbose        bool            `json:"verbose"`
	Wordlist       string          `json:"wordlist"`
	WordlistType   string          `json:"wordlist_type"`
}

type Fuzzer struct {
	Config    Config     `json:"config"`
	Request   Request    `json:"request"`
	Responses []Response `json:"response"`

	TotalWords int `json:"total_words"`
	Errors     int `json:"errors"`

	// mu *sync.Mutex `json:"-"`
}

// Initialize a new Fuzzer
func NewFuzzer(ctx context.Context, options cmd.CmdOptions, fuzztype string, wordlistType string, totalWords int) Fuzzer {
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
		HideWords:      options.HideWords,
		Host:           extractHost(options.URL),
		MatchLength:    options.MatchLength,
		MatchStatus:    options.MatchStatus,
		MatchWords:     options.MatchWords,
		Method:         options.Method,
		PostData:       options.PostData,
		Proxy:          options.Proxy,
		Retry:          options.Retry,
		Scan:           options.Scan,
		Threads:        options.Threads,
		Timeout:        options.Timeout,
		URL:            options.URL,
		UserAgent:      options.UserAgent,
		Verbose:        options.Verbose,
		Wordlist:       options.Wordlist,
		WordlistType:   wordlistType,
	}

	// Auto EGG
	if fuzztype == "" {
		f.Config.URL = util.AdjustUrlSuffix(f.Config.URL) + "EGG"
	}

	f.Request = NewRequest(f.Config)
	f.Responses = make([]Response, 0)
	f.TotalWords = totalWords
	f.Errors = 0
	return f
}

// Run to fuzz
func (f *Fuzzer) Run() error {
	runtime.GOMAXPROCS(f.Config.Threads)
	var wg sync.WaitGroup

	bar := *output.NewProgressBar(f.TotalWords, "Fuzzing...")

	wordCh := make(chan string, f.Config.Threads)

	output.Head("DIRECTORIES FOUND")

	// Wordlist from a file
	if f.Config.WordlistType == "" {

		readFile, err := os.Open(f.Config.Wordlist)
		if err != nil {
			log.Fatal(err)
		}
		defer readFile.Close()

		scanner := bufio.NewScanner(readFile)
		scanner.Split(bufio.ScanLines)

		for i := 0; i < f.Config.Threads; i++ {
			wg.Add(1)
			go f.worker(&wg, i, wordCh, bar)
		}

		for scanner.Scan() {
			bar.Add(1)
			bar.Describe(fmt.Sprintf("| Errors %d\r", f.Errors))

			word := scanner.Text()
			wordCh <- word

			// Auto EGG
			if f.Config.FuzzType == "" {
				// TXT files
				wordCh <- insertExt(word, ".txt")
				// HTML files
				wordCh <- insertExt(word, ".html")
				// PHP files
				wordCh <- insertExt(word, ".php")
				// Hidden files
				wordCh <- insertExt(word, ".")
			}

		}
	} else {
		// Built-in wordlist

		// Analyze
		words := strings.Split(f.Config.Wordlist, "_")
		wordArr := make([]string, 0)
		if words[0] == "ALPHA" {
			runes := []rune(words[1] + words[2])
			start := runes[0]
			end := runes[1]
			for i := start; i <= end; i++ {
				// Lowercase
				wordArr = append(wordArr, strings.ToLower(string(i)))
				// Uppercase
				wordArr = append(wordArr, strings.ToUpper(string(i)))

			}
		} else if words[0] == "NUM" {
			digits := len(words[1])

			start, err := strconv.Atoi(words[1])
			if err != nil {
				return fmt.Errorf("%s", err)
			}
			end, err := strconv.Atoi(words[2])
			if err != nil {
				return fmt.Errorf("%s", err)
			}

			// Create a numbers array
			for i := start; i <= end; i++ {
				wordArr = append(wordArr, fmt.Sprintf("%0*d", digits, i))
			}
		}

		// Reasign progressbar
		bar = *output.NewProgressBar(len(wordArr), "Fuzzing...")

		for i := 0; i < f.Config.Threads; i++ {
			wg.Add(1)
			go f.worker(&wg, i, wordCh, bar)
		}

		for _, w := range wordArr {
			bar.Add(1)
			bar.Describe(fmt.Sprintf("| Errors %d\r", f.Errors))
			wordCh <- w
		}
	}

	bar.Close()
	close(wordCh)
	wg.Wait()

	fmt.Println()

	// Scan contents
	if f.Config.Scan {
		if err := f.Scan(); err != nil {
			return err
		}
	}

	return nil
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
			f.printResult(resp, bar)
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
		if err != nil {
			f.Errors++
		} else {
			ok = true
		}
		time.Sleep(getDelay(f.Config.Delay))
		cnt++
	}
	return resp, nil
}

// Check if the response matches rules
func (f *Fuzzer) matcher(resp Response) bool {
	if util.ContainInt(f.Config.MatchStatus, resp.StatusCode) && !util.ContainInt(f.Config.HideStatus, resp.StatusCode) && f.matchContentLength(resp.ContentLength) && f.matchContentWords(resp.ContentWords) {
		return true
	}
	return false
}

// Print result
func (f *Fuzzer) printResult(resp Response, bar progressbar.ProgressBar) {
	bar.Clear()

	keyword := resp.Path

	if f.Config.FuzzType != "" && f.Config.FuzzType != "URL" {
		keyword = resp.Word
	}

	result := fmt.Sprintf("%-32s %s%s %-22s %s%s",
		color.CyanString(keyword),
		color.YellowString("("),
		color.GreenString("status: %d", resp.StatusCode),
		color.HiBlueString("length: %d", resp.ContentLength),
		color.HiMagentaString("words: %d", resp.ContentWords),
		color.YellowString(")"))

	if resp.RedirectPath != "" {
		fmt.Printf("%s -> %s\n", result, color.GreenString(resp.RedirectPath))
	} else {
		fmt.Printf("%s\n", result)
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

// Insert extentions (.php, .txt, etc.) to the word
// This is only used in Auto EGG mode.
func insertExt(word string, ext string) string {
	if ext == "." {
		// if the word is directory, insert "." before the last path.
		if strings.Contains(word, "/") {
			paths := strings.Split(word, "/")
			lastPath := "." + paths[len(paths)-1]
			return strings.Join(paths[:len(paths)-2], "/") + "/" + lastPath
		} else {
			return ext + word
		}
	} else {
		return word + ext
	}
}
