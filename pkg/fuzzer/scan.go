package fuzzer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/pkg/output"
)

func (f *Fuzzer) Scan() error {
	color.Yellow(output.TMPL_BAR_DOUBLE_M)
	fmt.Printf("%s SCAN WEB CONTENTS\n", color.YellowString("+"))
	color.Yellow(output.TMPL_BAR_DOUBLE_M)

	// keyword := "password"
	reKw := regexp.MustCompile(`(?i)admin|author|backup|cms|config|cred|pass|secret|sql|token|user|version|wordpress`)

	if len(f.Responses) > 0 {
		found := make(map[string][]string, 0)
		for _, resp := range f.Responses {
			foundLines := make([]string, 0)

			lines := strings.Split(resp.Content, "\n")
			for _, line := range lines {
				keyword := reKw.FindString(line)
				if keyword != "" {
					highlightLine := strings.TrimSpace(strings.ReplaceAll(line, keyword, color.HiCyanString(keyword)))
					foundLines = append(foundLines, highlightLine)
				}
			}

			if len(foundLines) > 0 {
				found[resp.Path] = foundLines
			}

		}

		// Output
		if len(found) > 0 {
			for key, val := range found {
				color.Cyan(key)
				for _, v := range val {
					fmt.Printf("%s %s\n", color.YellowString("|- "), v)
				}
			}
		}
	}
	return nil
}
