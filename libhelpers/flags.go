package libhelpers

import (
	"flag"
	"fmt"
	"net/url"
	"os"
)

const DEFAULT_WORDLIST = "/usr/share/seclists/Discovery/Web-Content/common.txt"

type Flags struct {
	TimeDelay string
	Url       string
	Verbose   bool
	Wordlist  string
}

func NewFlags(args []string) Flags {
	timeDelayPtr := flag.String("td", "1000-1500", "time delay per requests e.g. 1000ms. or random delay e.g. 1000ms-1500ms")
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

	var flags Flags
	flags.TimeDelay = *timeDelayPtr
	flags.Url = *urlPtr
	flags.Verbose = *verbosePtr
	flags.Wordlist = *wordlistPtr

	return flags
}
