package libhelpers

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

const DEFAULT_WORDLIST = "/usr/share/seclists/Discovery/Web-Content/common.txt"

func Flag(args []string) (string, bool, string) {
	urlPtr := flag.String("u", "", "target url")
	verbosePtr := flag.Bool("v", false, "verbose mode")
	wordlistPtr := flag.String("w", DEFAULT_WORDLIST, "the wordlist to use")

	flag.Parse()

	if *urlPtr == "" {
		fmt.Printf("URL is not specified.\n\n")
		flag.CommandLine.Usage()
		os.Exit(0)
	} else {
		u, _ := url.ParseRequestURI(*urlPtr)
		if u == nil {
			fmt.Printf("The specified url is invalid.\n\n")
			flag.CommandLine.Usage()
			os.Exit(0)
		}
	}

	return *urlPtr, *verbosePtr, *wordlistPtr
}
