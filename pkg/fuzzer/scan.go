package fuzzer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/pkg/output"
	"github.com/hideckies/fuzzagotchi/pkg/util"
)

func (f *Fuzzer) Scan() error {
	if !f.Config.Scan {
		return nil
	}

	reKwHeader := regexp.MustCompile(`(?i)apache|nginx|php|python|werkzeug|\d+\.\d+\.\d+`)
	reKwContent := regexp.MustCompile(`(?i)admin|author|backup|cms|config|cred|pass|root|secret|sql|token|user|version|wordpress|\d+\.\d+\.\d+`)

	if len(f.Responses) > 0 {
		foundHeaders := make(map[string]string, 0)
		foundContents := make(map[string][]string, 0)

		// To prevent contents duplicate
		uniqContentLines := make([]string, 0)

		for _, resp := range f.Responses {
			// Headers
			for key, header := range resp.Header {
				if key == "Accept-Ranges" || key == "Cache-Control" || key == "Connection" || key == "Content-Length" || key == "Content-Type" || key == "Date" || key == "Etag" || key == "Expires" || key == "Last-Modified" || key == "Location" || key == "Keep-Alive" || key == "Pragma" || key == "Vary" {
					continue
				}
				val := strings.Join(header, " ")
				keyword := reKwHeader.FindString(val)
				if _, ok := foundHeaders[key]; !ok {
					foundHeaders[key] = strings.ReplaceAll(val, keyword, color.YellowString(keyword))
				}
			}

			// Contents
			tmpContentLines := make([]string, 0)
			lines := strings.Split(resp.Content, "\n")

			for _, line := range lines {
				keyword := reKwContent.FindString(line)
				if keyword != "" && !util.ContainString(uniqContentLines, line) {
					highlightLine := strings.TrimSpace(strings.ReplaceAll(line, keyword, color.YellowString(keyword)))
					tmpContentLines = append(tmpContentLines, highlightLine)
					uniqContentLines = append(uniqContentLines, line)
				}
			}

			if len(tmpContentLines) > 0 {
				foundContents[resp.Path] = tmpContentLines
			}

		}

		// Output header
		output.Head("RESPONSE HEADERS")
		if len(foundHeaders) > 0 {
			for key, val := range foundHeaders {
				fmt.Printf("%s: %s\n", key, val)
			}
		} else {
			fmt.Println("Nothing found.")
		}

		// Output contents
		output.Head("SENSITIVE DATA IN CONTENTS")
		if len(foundContents) > 0 {
			for key, val := range foundContents {
				color.Cyan(key)
				for _, v := range val {
					fmt.Printf("%s %s\n", color.YellowString("|- "), v)
				}
			}
		} else {
			fmt.Println("Nothing found.")
		}
	}
	return nil
}
