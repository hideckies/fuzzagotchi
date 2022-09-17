package libfuzz

import (
	"bufio"
	"fmt"
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

// Fuzz fuzzes on the content discovery.
func Fuzz(conf libgotchi.Conf) {
	var wg sync.WaitGroup

	readFile, err := os.Open(conf.Wordlist)
	if err != nil {
		color.HiRed("%v\nPlease install seclists by running 'sudo apt install seclists'.\n", err)
		os.Exit(0)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	resCh := make(chan libgotchi.Res)

	words := make([]string, 0)

	for fileScanner.Scan() {
		wg.Add(1)
		word := fileScanner.Text()
		words = append(words, word)
		go Process(&wg, conf, word, resCh)
	}

	readFile.Close()

	for i := 0; i < len(words); i++ {
		res := <-resCh
		Output(res)
	}
	wg.Wait()
}

func Process(wg *sync.WaitGroup, conf libgotchi.Conf, word string, resCh chan libgotchi.Res) {
	defer wg.Done()

	req := libgotchi.NewReq(conf)

	time.Sleep(req.Duration)

	if len(conf.Header) > 0 {
		headers := strings.Split(conf.Header, ";")
		for _, v := range headers {
			header := strings.Split(strings.TrimSpace(v), ":")
			key := header[0]
			val := header[1]
			req.Headers[key] = val
		}
	}
	// Update cookies
	if len(conf.Cookie) > 0 {
		cookies := strings.Split(conf.Cookie, ";")
		for _, v := range cookies {
			c := strings.Split(strings.TrimSpace(v), "=")
			key := c[0]
			val := c[1]
			req.Cookies[key] = val
		}
	}

	// Send request
	res := req.Send(word)
	resCh <- res
}

// Print results
func Output(res libgotchi.Res) {
	result := fmt.Sprintf(
		"%-10s\t\tStatus: %d, Content Length: %d, Duration: %.2fs",
		res.Word,
		res.StatusCode,
		res.ContentLength,
		res.Duration.Abs().Seconds())
	resultFailed := fmt.Sprintf("[x] %v", result)
	if res.Config.Color {
		result = color.HiGreenString(result)
		resultFailed = color.RedString(resultFailed)
	}

	rcl, _ := regexp.Compile("^([1-9][0-9]*|0)$")
	rclrange, _ := regexp.Compile("^(([1-9][0-9]*|0)-([1-9][0-9]*|0))$")
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
