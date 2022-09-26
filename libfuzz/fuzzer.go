package libfuzz

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/libgotchi"
	"github.com/hideckies/fuzzagotchi/libutils"
)

type Fuzzer struct {
	Config libgotchi.Conf
}

// Run executes a fuzzing
func (f *Fuzzer) Run() {
	var wg sync.WaitGroup
	wg.Add(f.Config.Threads)

	readFile, err := os.Open(f.Config.Wordlist)
	if err != nil {
		log.Fatal(err)
		return
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

func (f *Fuzzer) worker(wg *sync.WaitGroup, id int, wordCh chan string) {
	defer wg.Done()
	for {
		select {
		case word, ok := <-wordCh:
			if !ok {
				return
			}
			f.process(word)
			select {
			case <-time.After(libgotchi.NewRate(f.Config.Rate)):
			}
		}
	}
}

func (f *Fuzzer) process(word string) {
	req := libgotchi.NewReq(f.Config)

	time.Sleep(req.Rate)

	// Send request
	res, err := req.Send(word)
	if err != nil {
		if f.Config.Verbose {
			fmt.Printf("%-10s\t\tError: %s\n", word, err)
		}
		f.process(word)
		return
	}

	f.output(res)
}

// Print results
func (f *Fuzzer) output(res libgotchi.Res) {
	result := fmt.Sprintf(
		"%-10s\t\tStatus: %d, Content Length: %d, Duration: %.2fs",
		res.Word,
		res.StatusCode,
		res.ContentLength,
		res.Rate.Abs().Seconds())
	resultFailed := fmt.Sprintf("[x] %v", result)
	if res.Config.Color {
		result = color.HiGreenString(result)
		resultFailed = color.RedString(resultFailed)
	}

	rcl, _ := regexp.Compile("^([1-9][0-9]*|0)$")
	rclrange, _ := regexp.Compile("^(([1-9][0-9]*|0)-([1-9][0-9]*|0))$")

	if rcl.MatchString(res.Config.NoContentLength) {
		ncl, _ := strconv.Atoi(res.Config.NoContentLength)
		if libutils.IntContains(res.Config.Status, res.StatusCode) && ncl != res.ContentLength {
			fmt.Println(result)
		} else if res.Config.Verbose {
			fmt.Println(resultFailed)
		}
	} else {
		if rclrange.MatchString(res.Config.ContentLength) {
			contentlengths := strings.Split(res.Config.ContentLength, "-")
			cmin, _ := strconv.Atoi(contentlengths[0])
			cmax, _ := strconv.Atoi(contentlengths[1])
			if libutils.IntContains(res.Config.Status, res.StatusCode) && (cmin <= res.ContentLength && res.ContentLength <= cmax) {
				fmt.Println(result)
			} else if res.Config.Verbose {
				fmt.Println(resultFailed)
			}
		} else if rcl.MatchString(res.Config.ContentLength) {
			cl, _ := strconv.Atoi(res.Config.ContentLength)
			if libutils.IntContains(res.Config.Status, res.StatusCode) && cl == res.ContentLength {
				fmt.Println(result)
			} else if res.Config.Verbose {
				fmt.Println(resultFailed)
			}
		} else {
			if libutils.IntContains(res.Config.Status, res.StatusCode) {
				fmt.Println(result)
			} else if res.Config.Verbose {
				fmt.Println(resultFailed)
			}
		}
	}
}

// NewFuzzer returns a new Fuzzer
func NewFuzzer(conf libgotchi.Conf) Fuzzer {
	var f Fuzzer
	f.Config = conf
	return f
}
