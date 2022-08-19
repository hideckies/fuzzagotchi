package libdir

import (
	"bufio"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/libgotchi"
	"github.com/hideckies/fuzzagotchi/libhelpers"
)

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

	// Time delay
	duration := libgotchi.NewDuration(flags)

	var word string
	for fileScanner.Scan() {
		time.Sleep(duration)

		word = fileScanner.Text()

		// ******************************************************************
		// Request
		// Update the request configuration
		reqConfPtr.Url = strings.Replace(flags.Url, "EGG", word, -1)
		// Send request
		res := libgotchi.SendRequest(reqConfPtr)
		// ******************************************************************

		// Display result
		switch res.StatusCode {
		case 200:
			color.HiGreen("%-30s\t\tStatus Code: %d, Content Length: %d", word, res.StatusCode, res.ContentLength)
		case 301, 302:
			color.HiGreen("%-30s\t\tStatus Code: %d, Content Length: %d", word, res.StatusCode, res.ContentLength)
		case 400, 401, 402, 403, 404, 405:
			if flags.Verbose {
				color.Red("%-30s\t\tStatus Code: %d, Content Length: %d", word, res.StatusCode, res.ContentLength)
			}
		case 500:
			if flags.Verbose {
				color.Red("%-30s\t\tStatus Code: %d, Content Length: %d", word, res.StatusCode, res.ContentLength)
			}
		default:
			if flags.Verbose {
				color.White("%-30s\t\tStatus Code: %d, Content Length: %d", word, res.StatusCode, res.ContentLength)
			}
		}
	}

	readFile.Close()
}
