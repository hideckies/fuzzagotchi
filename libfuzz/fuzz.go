package libfuzz

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/libgotchi"
	"github.com/hideckies/fuzzagotchi/libhelpers"
	"github.com/hideckies/fuzzagotchi/libutils"
)

// Fuzz fuzzes on the content discovery.
func Fuzz(flags libhelpers.Flags) {
	readFile, err := os.Open(flags.Wordlist)
	if err != nil {
		color.HiRed("%v\nPlease install seclists by running 'sudo apt install seclists'.\n", err)
		os.Exit(0)
	}
	fileScanner := bufio.NewScanner(readFile)
	fileScanner.Split(bufio.ScanLines)

	// Initialize a configuration for request
	reqConf := libgotchi.NewReqConf()
	reqConfPtr := &reqConf

	var word string
	for fileScanner.Scan() {
		duration := libgotchi.NewDuration(flags)
		time.Sleep(duration)

		word = fileScanner.Text()

		// ******************************************************************
		// Request
		// ******************************************************************

		// Update url
		reqConfPtr.Url = flags.Url
		// Update method
		reqConfPtr.Method = flags.Method
		// Update headers
		if len(flags.Header) > 0 {
			headers := strings.Split(flags.Header, ";")
			for _, v := range headers {
				header := strings.Split(strings.TrimSpace(v), ":")
				key := header[0]
				val := header[1]
				reqConfPtr.Headers[key] = val
			}
		}
		// Update cookies
		if len(flags.Cookie) > 0 {
			cookies := strings.Split(flags.Cookie, ";")
			for _, v := range cookies {
				c := strings.Split(strings.TrimSpace(v), "=")
				key := c[0]
				val := c[1]
				reqConfPtr.Cookies[key] = val
			}
		}

		// Send request
		res := libgotchi.SendRequest(reqConfPtr, word)
		// ******************************************************************

		result := fmt.Sprintf(
			"%-40s\t\tStatus: %d, Content Length: %d, Duration: %.2fs",
			word,
			res.StatusCode,
			res.ContentLength,
			duration.Abs().Seconds())
		resultFailed := fmt.Sprintf("[x] %v", result)
		if flags.Color {
			result = color.HiGreenString(result)
			resultFailed = color.RedString(resultFailed)
		}

		// ******************************************************************
		// Display the result
		// Filter the content length & verbose
		// ******************************************************************
		rcl, _ := regexp.Compile("^([1-9][0-9]*|0)$")
		rclrange, _ := regexp.Compile("^(([1-9][0-9]*|0)-([1-9][0-9]*|0))$")
		if rclrange.MatchString(flags.ContentLength) {
			contentlengths := strings.Split(flags.ContentLength, "-")
			cmin, _ := strconv.Atoi(contentlengths[0])
			cmax, _ := strconv.Atoi(contentlengths[1])
			if libutils.IntContains(flags.Status, res.StatusCode) && (cmin <= res.ContentLength && res.ContentLength <= cmax) {
				fmt.Println(result)
			} else if flags.Verbose {
				fmt.Println(resultFailed)
			}
		} else if rcl.MatchString(flags.ContentLength) {
			cl, _ := strconv.Atoi(flags.ContentLength)
			if libutils.IntContains(flags.Status, res.StatusCode) && cl == res.ContentLength {
				fmt.Println(result)
			} else if flags.Verbose {
				fmt.Println(resultFailed)
			}
		} else {
			if libutils.IntContains(flags.Status, res.StatusCode) {
				fmt.Println(result)
			} else if flags.Verbose {
				fmt.Println(resultFailed)
			}
		}
		// ******************************************************************
	}

	readFile.Close()
}
