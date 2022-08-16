package main

import (
	"os"

	"github.com/fatih/color"

	"github.com/hideckies/fuzzagotchi/libdir"
	"github.com/hideckies/fuzzagotchi/libhelpers"
)

func main() {
	flags := libhelpers.NewFlags(os.Args)

	color.HiCyan("%s\n\n", libhelpers.LOGO)
	color.HiCyan("Target URL: %s\n", flags.Url)
	color.HiCyan("Wordlist: %s\n", flags.Wordlist)
	color.HiCyan("Verbose: %t\n", flags.Verbose)
	color.HiCyan("%s\n\n", libhelpers.BAR_DOUBLE_M)

	libdir.Fuzz(flags)
}
