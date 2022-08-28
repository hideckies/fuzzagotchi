package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/libfuzz"
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
	fuzzagotchi -u https://example.com/EGG -w wordlist.txt --status 200,301
	fuzzagotchi -u https://example.com/EGG -w wordlist.txt --content-length 120-175
	fuzzagotchi -u https://example.com/EGG -w wordlist.txt -H "Authorization: Bearer <token>"
	fuzzagotchi -u https://example.com/EGG -w wordlist.txt -C "name1=value1; name2=value2"
	
	fuzzagotchi -u https://example.com/EGG.php -w wordlist.txt
	fuzzagotchi -u https://example.com/?q=EGG -w wordlist.txt

  [Subdomain Scan] *Under development so unavailable currently.
	fuzzagotchi -u https://example.com -w subdomains.txt -H "Host: EGG.example.com" --content-length 500-2000

  [Brute Force Credentials] *Unser development so unavailable currently.
	fuzzagotchi -M POST -u https://example.com/login -w passwords.txt --post-data "username=admin&password=EGG"
	fuzzagotchi -M POST -u https://example.com/login -w passwords.txt --post-data "{username:admin, password: EGG}"
	`,
}

const DEFAULT_WORDLIST = "/usr/share/seclists/Discovery/Web-Content/common.txt"

func main() {
	flags := libhelpers.NewFlags()

	cmd.Flags().BoolVarP(&flags.Color, "color", "", false, "The output color")
	cmd.Flags().StringVarP(&flags.ContentLength, "content-length", "", "-1", "Display the specific content length e.g. 120-560")
	cmd.Flags().StringVarP(&flags.Cookie, "cookie", "C", "", "Custom cookie e.g. \"name1=value1; name2=value2\"")
	cmd.Flags().StringVarP(&flags.TimeDelay, "delay", "", "0.2-0.5", "Time delay (seconds) per requests e.g. 1.2. Or random delay e.g. 0.8-1.5")
	cmd.Flags().StringVarP(&flags.Header, "header", "H", "", "Custom header e.g. \"Authorization: Bearer <token>; Host: example.com\"")
	cmd.Flags().StringVarP(&flags.Method, "method", "M", "GET", "Specific method e.g. GET, POST, PUT, OPTIONS, etc.")
	cmd.Flags().StringVarP(&flags.PostData, "post-data", "", "", "POST request with data e.g. \"username=admin&password=EGG\"")
	cmd.Flags().IntSliceVarP(&flags.Status, "status", "s", []int{200, 204, 301, 302, 307, 401, 403}, "Display the specific status codes only")
	cmd.Flags().Int8VarP(&flags.Threads, "threads", "t", 10, "Number of concurrent threads.")
	cmd.Flags().StringVarP(&flags.Url, "url", "u", "", "Target URL (required)")
	cmd.MarkFlagRequired("url")
	cmd.Flags().BoolVarP(&flags.Verbose, "verbose", "v", false, "Verbose mode")
	cmd.Flags().StringVarP(&flags.Wordlist, "wordlist", "w", DEFAULT_WORDLIST, "Wordlist for fuzzing")

	cmd.Run = func(cmd1 *cobra.Command, args []string) {
		var s []string
		s = append(s, fmt.Sprintf("%s\n\n", libhelpers.LOGO))
		s = append(s, fmt.Sprintf("%-10s\t\t%t\n", "Output Color:", flags.Color))
		s = append(s, fmt.Sprintf("%-10s\t\t%s\n", "URL:", flags.Url))
		s = append(s, fmt.Sprintf("%-10s\t\t%s\n", "Wordlist:", flags.Wordlist))
		s = append(s, fmt.Sprintf("%-10s\t\t%s\n", "Method:", flags.Method))
		if len(flags.Header) > 0 {
			s = append(s, fmt.Sprintf("%-10s\t\t%s\n", "Header:", flags.Header))
		}
		if len(flags.Cookie) > 0 {
			s = append(s, fmt.Sprintf("%-10s\t\t%s\n", "Cookie:", flags.Cookie))
		}
		s = append(s, fmt.Sprintf("%-10s\t\t%v\n", "Status:", strings.Trim(strings.Replace(fmt.Sprint(flags.Status), " ", ",", -1), "[]")))
		cl, _ := strconv.Atoi(flags.ContentLength)
		if cl >= 0 {
			s = append(s, fmt.Sprintf("%-10s\t\t%s\n", "Content Length:", flags.ContentLength))
		}
		s = append(s, fmt.Sprintf("%-10s\t\t%s\n", "Duration: ", flags.TimeDelay))
		s = append(s, fmt.Sprintf("%-10s\t\t%t\n", "Verbose:", flags.Verbose))
		s = append(s, fmt.Sprintf("%s\n\n", libhelpers.BAR_DOUBLE_L))

		for _, v := range s {
			if flags.Color {
				fmt.Printf(color.HiCyanString(v))
			} else {
				fmt.Printf(v)
			}
		}

		libfuzz.Fuzz(flags)
	}

	if err := cmd.Execute(); err != nil {
		os.Exit(0)
	}
}
