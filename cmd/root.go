package cmd

import (
	"github.com/spf13/cobra"
)

var (
	Proceed bool

	rootCmd = &cobra.Command{
		Use: "fuzzagotchi",
		Run: func(cmd *cobra.Command, args []string) {
			Proceed = true
		},
	}
)

const USAGE_TEMPLATE = `FUZZAGOTCHI - Automatic Web Fuzzer

Fuzzagotchi is so methodical and looks for details without you asking.
/EGG, /.EGG, /EGG.txt, /EGG.html, /EGG.php, etc. This tool is automatic and exhaustive.


USAGE:
  fuzzagotchi -u <URL> -w <WORDLIST> [OPTIONS]

  -u, --url              URL to fuzz
  -w, --wordlist         Wordlist used for fuzzing

FUZZER OPTIONS:
  -X, --method           HTTP method (default: GET)
  -H, --header           Custom header
  -C, --cookie           Custom cookie
  -d, --post-data        POST data
  -f, --follow-redirect  Follow redirects (default: false)
  -p, --delay            Delay between each request e.g. 0.8-1.5
      --retry            Number of retry when a request error (default: 2)
  -r, --recursion        Enable resursive brute force (default: false)
  -t, --threads          Number of threads (default: 20)
      --timeout          Request timeout in seconds (default: 10)
      --user-agent       Specific User-Agent
  -x, --proxy            Proxy URL e.g. http://127.0.0.1:8080
  
      --ms               Display only matched status code (default: 200, 204, 301, 302, 307, 401, 403, 500)
      --ml               Display only matched content-length
      --mw               Display only matched the number of words
	  --mk               Display only the response content contains specific keyword
      --hs               Hide matched status code
      --hl               Hide matched content-length
      --hw               Hide matched the number of words
	  --hk               Hide if the response content contains specific keyword
  
      --no-color         Disable colorize the output (default: false)
      -v, --verbose      Verbose mode (default: false)

Deep Options:
  -s, --scan             Scan sensitive information, vulnerabilities.

META OPTIONS:
  -h, --help             Print the usage of Fuzzagotchi
  version                Print the version of Fuzzagotchi

BUILTIN WORDLISTS:
  ALPHA_A_Z              Alphabets from a to z (contains both lowercase and uppercase)
  NUM_0000_9999          Numbers from min to max

EXAMPLES:
  fuzzagotchi -u https://example.com -w wordlist.txt
  fuzzagotchi -u https://example.com/EGG -w wordlist.txt
  fuzzagotchi -u https://example.com/?id=EGG -w NUM_0_999
  fuzzagotchi -u https://example.com -w wordlist.txt --scan
`

type CmdOptions struct {
	URL      string
	Wordlist string

	Method   string
	Header   string
	Cookie   string
	PostData string

	FollowRedirect bool
	Delay          string
	Recursion      bool
	Retry          int
	Threads        int
	Timeout        int
	UserAgent      string

	MatchStatus  string
	MatchLength  string
	MatchWords   string
	MatchKeyword string
	HideStatus   string
	HideLength   string
	HideWords    string
	HideKeyword  string

	Proxy string

	Scan bool

	NoColor bool
	Verbose bool
}

var Options = CmdOptions{}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.Flags().StringVarP(&Options.URL, "url", "u", "", "")
	rootCmd.MarkFlagRequired("url")
	rootCmd.Flags().StringVarP(&Options.Wordlist, "wordlist", "w", "", "")
	rootCmd.MarkFlagRequired("wordlist")

	rootCmd.Flags().StringVarP(&Options.Method, "method", "X", "GET", "")
	rootCmd.Flags().StringVarP(&Options.Header, "header", "H", "", "")
	rootCmd.Flags().StringVarP(&Options.Cookie, "cookie", "C", "", "")
	rootCmd.Flags().StringVarP(&Options.PostData, "post-data", "d", "", "")

	rootCmd.Flags().BoolVarP(&Options.FollowRedirect, "follow-redirect", "f", false, "")
	rootCmd.Flags().StringVarP(&Options.Delay, "delay", "p", "0", "")
	rootCmd.Flags().BoolVarP(&Options.Recursion, "recursion", "r", false, "")
	rootCmd.Flags().IntVarP(&Options.Retry, "retry", "", 2, "")
	rootCmd.Flags().IntVarP(&Options.Threads, "threads", "t", 20, "")
	rootCmd.Flags().IntVarP(&Options.Timeout, "timeout", "", 10, "")
	rootCmd.Flags().StringVarP(&Options.UserAgent, "user-agent", "", "Fuzzagotchi", "")

	rootCmd.Flags().StringVarP(&Options.MatchStatus, "ms", "", "200,204,301,302,307,401,403,500", "")
	rootCmd.Flags().StringVarP(&Options.MatchLength, "ml", "", "", "")
	rootCmd.Flags().StringVarP(&Options.MatchWords, "mw", "", "", "")
	rootCmd.Flags().StringVarP(&Options.MatchKeyword, "mk", "", "", "")
	rootCmd.Flags().StringVarP(&Options.HideStatus, "hs", "", "", "")
	rootCmd.Flags().StringVarP(&Options.HideLength, "hl", "", "", "")
	rootCmd.Flags().StringVarP(&Options.HideWords, "hw", "", "", "")
	rootCmd.Flags().StringVarP(&Options.HideKeyword, "hk", "", "", "")

	rootCmd.Flags().StringVarP(&Options.Proxy, "proxy", "x", "", "")

	rootCmd.Flags().BoolVarP(&Options.Scan, "scan", "s", false, "")

	rootCmd.Flags().BoolVarP(&Options.NoColor, "no-color", "", false, "")
	rootCmd.Flags().BoolVarP(&Options.Verbose, "verbose", "v", false, "")

	// Set custom usage
	rootCmd.SetUsageTemplate(USAGE_TEMPLATE)
}
