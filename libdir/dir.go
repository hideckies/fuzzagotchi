package libdir

import (
	"bufio"
	"os"
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

	var word string
	for fileScanner.Scan() {
		// Time delay
		duration := libgotchi.NewDuration(flags)
		time.Sleep(duration)

		word = fileScanner.Text()

		// Request
		// Create a configuration for request
		reqConf := libgotchi.NewReqConf()
		reqConfPtr := &reqConf
		// Update the request configuration
		reqConfPtr.Url = libgotchi.AdjustUrlSuffix(flags.Url) + word
		// Send request
		res := libgotchi.SendRequest(reqConfPtr)

		// Display result
		switch res.StatusCode {
		case 200, 302:
			color.HiGreen("%s: %s", word, res.Status)
		case 400, 401, 402, 403, 404, 405:
			if flags.Verbose {
				color.Red("%s: %s", word, res.Status)
			}
		case 500:
			if flags.Verbose {
				color.Red("%s: %s", word, res.Status)
			}
		default:
			if flags.Verbose {
				color.White("%s: %s", word, res.Status)
			}
		}
	}

	readFile.Close()
}
