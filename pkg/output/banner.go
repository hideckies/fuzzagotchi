package output

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/cmd"
)

const (
	banner = `
 ____  _  _  ____  ____   __    ___   __  ____  ___  _  _  __  
(  __)/ )( \(__  )(__  ) / _\  / __) /  \(_  _)/ __)/ )( \(  ) 
 ) _) ) \/ ( / _/  / _/ /    \( (_ \(  O ) )( ( (__ ) __ ( )(  
(__)  \____/(____)(____)\_/\_/ \___/ \__/ (__) \___)\_)(_/(__) 
`

	subtitle = `
	                                Automatic Web Fuzzer
`
)

func Banner(options cmd.CmdOptions) {
	color.Yellow(banner)
	color.Cyan(subtitle)

	fmt.Println()

	plusMark := color.CyanString("+")

	fmt.Printf("%s%s%s\n", plusMark, color.YellowString(TMPL_BAR_SINGLE_M), plusMark)
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', tabwriter.TabIndent)
	fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("URL"), color.CyanString(options.URL))
	fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("Wordlist"), color.CyanString(options.Wordlist))
	fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("Method"), color.CyanString(options.Method))
	if len(options.Header) > 0 {
		fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("Header"), color.CyanString(options.Header))
	}
	fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("Threads"), color.CyanString("%d", options.Threads))
	fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("Delay"), color.CyanString(options.Delay))
	fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("Follow Redirect"), color.CyanString("%t", options.FollowRedirect))
	fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("Recursion"), color.CyanString("%t", options.Recursion))
	fmt.Fprintf(w, "%s %s\t:\t%s\n", plusMark, color.YellowString("Verbose"), color.CyanString("%t", options.Verbose))
	w.Flush()
	fmt.Printf("%s%s%s\n", plusMark, color.YellowString(TMPL_BAR_SINGLE_M), plusMark)
	fmt.Println()
}
