package libdir

import (
	"bufio"
	"fmt"
	"os"
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

		// Display results
		resultSuccess := color.HiGreenString("%-30s\t\tStatus Code: %d, Content Length: %d", word, res.StatusCode, res.ContentLength)
		resultFailed := color.RedString("%-30s\t\tStatus Code: %d, Content Length: %d", word, res.StatusCode, res.ContentLength)
		if flags.ContentLength > 0 && flags.ContentLength == res.ContentLength {
			fmt.Println(resultSuccess)
		} else if flags.ContentLength == -1 && libutils.IntContains(flags.StatusCodes, res.StatusCode) {
			fmt.Println(resultSuccess)
		} else if flags.Verbose {
			fmt.Println(resultFailed)
		}
	}

	readFile.Close()
}
