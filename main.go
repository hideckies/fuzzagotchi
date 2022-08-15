package main

import (
	"os"

	"github.com/fatih/color"

	"github.com/hideckies/fuzzagotchi/libdir"
	"github.com/hideckies/fuzzagotchi/libhelpers"
)

func main() {
	url, verbose, wordlist := libhelpers.Flag(os.Args)

	// Display ascii art
	color.HiCyan("%s\n\n", libhelpers.LOGO)

	color.HiCyan("Target URL: %s\n", url)
	color.HiCyan("Wordlist: %s\n", wordlist)
	color.HiCyan("%s\n\n", libhelpers.BAR_DOUBLE_M)

	// Start fuzzing
	libdir.Fuzz(url, verbose, wordlist)
}
