package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/hideckies/fuzzagotchi/pkg/bars"
	"github.com/hideckies/fuzzagotchi/pkg/flags"
	"github.com/hideckies/fuzzagotchi/pkg/fuzz"
)

var ascii = `
|~~|   |~~/~~/  /\   /~~\ /~~\ ~~|~~ /~~|  |~|~
|--|   | /  /  /__\ |  __|    |  |  |   |--| | 
|   \_/ /__/__/    \ \__/ \__/   |   \__|  |_|_
`

func main() {
	url, verbose, wordlist := flags.Flag(os.Args)

	// Display ascii art
	color.HiCyan(ascii)
	fmt.Println()

	color.HiCyan("Target URL: %s\n", url)
	color.HiCyan("Wordlist: %s\n", wordlist)
	color.HiCyan(bars.BAR_DOUBLE_M)
	fmt.Println()

	// Start fuzzing
	fuzz.Fuzz(url, verbose, wordlist)
}
