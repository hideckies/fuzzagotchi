package main

import (
	"os"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/libdir"
	"github.com/hideckies/fuzzagotchi/libhelpers"
	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:     "fuzzagotchi",
	Version: "v0.1.0",
	Short:   "A fuzzing tool written in Go.",
	Example: `
  [Content Discovery]
  fuzzagotchi -u https://example.com/EGG -w wordlist.txt
  fuzzagotchi -u https://example.com/EGG -w wordlist.txt -H "Cookie: isGotchi=true"
  fuzzagotchi -u https://example.com/EGG.php -w wordlist.txt
  fuzzagotchi -u https://example.com/?q=EGG -w wordlist.txt

  [Brute Force POST Data] *Unser development so unavailable currently.
  fuzzagotchi -u https://example.com/login -w passwords.txt --post-data "username=admin&password=EGG"
  fuzzagotchi -u https://example.com/login -w passwords.txt --post-data "{username:admin, password: EGG}"
  
  [Subdomain Scan] *Under development so unavailable currently.
  fuzzagotchi -u https://EGG.example.com -w wordlist.txt
	`,
}

const DEFAULT_WORDLIST = "/usr/share/seclists/Discovery/Web-Content/common.txt"

func main() {
	flags := libhelpers.NewFlags()

	cmd.Flags().StringVarP(&flags.TimeDelay, "delay", "", "100-200", "Time delay per requests e.g. 500ms. Or random delay e.g. 500ms-700ms")
	cmd.Flags().StringVarP(&flags.Header, "header", "H", "", "Custom header e.g. \"Authorization: Bearer <token>; Cookie: key=value\"")
	cmd.Flags().Int8VarP(&flags.Threads, "threads", "t", 10, "Number of concurrent threads")
	cmd.Flags().StringVarP(&flags.PostData, "post-data", "", "", "POST request with data e.g. \"username=admin&password=EGG\"")
	cmd.Flags().StringVarP(&flags.Url, "url", "u", "", "Target URL (required)")
	cmd.MarkFlagRequired("url")
	cmd.Flags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose mode")
	cmd.Flags().StringVarP(&flags.Wordlist, "wordlist", "w", DEFAULT_WORDLIST, "Wordlist for fuzzing")

	cmd.Run = func(cmd1 *cobra.Command, args []string) {
		color.HiCyan("%s\n\n", libhelpers.LOGO)
		color.HiCyan("Target URL: %s\n", flags.Url)
		color.HiCyan("Wordlist: %s\n", flags.Wordlist)
		color.HiCyan("Verbose: %t\n", flags.Verbose)
		color.HiCyan("%s\n\n", libhelpers.BAR_DOUBLE_L)

		libdir.Fuzz(flags)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}
