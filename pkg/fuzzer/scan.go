package fuzzer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/pkg/output"
	"github.com/hideckies/fuzzagotchi/pkg/util"
)

// Scan each page
func (f *Fuzzer) Scan() error {
	if !f.Config.Scan {
		return nil
	}

	reKwHeader := regexp.MustCompile(`
		(?i)apache|gunicorn|nginx|php|python|tomcat|werkzeug|wsgiserver|\d+\.\d+\.\d+`)
	reKwContent := regexp.MustCompile(`
		(?i)admin|author|backup|cms|config|cred|login|ng-apps|pass|root|secret|sql|token|user|version|wordpress|\d+\.\d+(\.\d+)*`)

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

		// Output header & vulnerabilities
		output.Head("RESPONSE HEADERS")
		if len(foundHeaders) > 0 {
			for key, val := range foundHeaders {
				fmt.Printf("%s: %s\n", key, val)
				// check if vulnerabilities
				vulns := f.FindVulns(val)
				if len(vulns) > 0 {
					for _, vuln := range vulns {
						fmt.Printf("%s\t%s %s\n", color.YellowString("|_"), color.HiMagentaString("->"), color.HiMagentaString(vuln))
					}
				}
			}
		} else {
			fmt.Println("Nothing found.")
		}

		// Output contents
		output.Head("WEB CONTENTS")
		if len(foundContents) > 0 {
			for key, val := range foundContents {
				color.Cyan(key)
				for _, v := range val {
					fmt.Printf("%s %s\n", color.YellowString("|_ "), v)
					// check if vulnerabilities
					vulns := f.FindVulns(v)
					if len(vulns) > 0 {
						for _, vuln := range vulns {
							fmt.Printf("%s\t%s %s\n", color.YellowString("|_"), color.HiMagentaString("->"), color.HiMagentaString(vuln))
						}
					}
				}
			}
		} else {
			fmt.Println("Nothing found.")
		}
	}
	return nil
}

// Find vulnerabilities
func (f *Fuzzer) FindVulns(val string) []string {
	vulns := make([]string, 0)

	reV := regexp.MustCompile(`\d+\.\d+\.\d+`)

	// Angular
	if strings.Contains(val, "ng-apps") {
		s := "Angular Server-Side Template Injection"
		vulns = append(vulns, s)
	}

	// Apache
	if strings.Contains(val, "Apache") {
		mVer := reV.FindString(val)
		if mVer == "2.4.49" {
			s := "Apache 2.4.49 Path Traversal (CVE-2021-41773)"
			vulns = append(vulns, s)
		}
		if mVer == "2.4.50" {
			s := "Apache 2.4.50 Path Traversal (CVE-2021-42013)"
			vulns = append(vulns, s)
		}

		// Tomcat
		if strings.Contains(val, "Tomcat") {
			s1 := "Apache Tomcat Remote Code Execution"
			vulns = append(vulns, s1)
			s2 := "Apache Tomcat AJP 'Ghostcat' File Inclusion"
			vulns = append(vulns, s2)
		}
	}

	// Exiftool
	if strings.Contains(val, "Exiftool") {
		s := "Exiftool Command Injection version < 12.38"
		vulns = append(vulns, s)
	}

	// Flask
	if strings.Contains(val, "Flask") {
		s := "Flask Server-Side Template Injection"
		vulns = append(vulns, s)
	}

	// Spring
	if strings.Contains(val, "Whitelabel") {
		s1 := "Spring Boot Server-Side Template Injection"
		vulns = append(vulns, s1)
		s2 := "Spring4Shell (CVE-2022-22965)"
		vulns = append(vulns, s2)
	}

	// Werkzeug
	if strings.Contains(val, "Werkzeug") {
		s := "Werkzeug Console Remote Code Execution"
		vulns = append(vulns, s)
	}

	// WordPress
	if strings.Contains(val, "WordPress") {
		mVer := reV.FindString(val)
		if strings.Contains(mVer, "5.6") || strings.Contains(mVer, "5.7") {
			s := "WordPress Authenticated XXE (CVE-2021-29447)"
			vulns = append(vulns, s)
		}
	}

	// Exclude existing vulnerabilities
	uniqVulns := make([]string, 0)
	for _, vuln := range vulns {
		if !util.ContainString(f.Vulns, vuln) {
			uniqVulns = append(uniqVulns, vuln)
		}
	}

	f.Vulns = append(f.Vulns, uniqVulns...)

	return uniqVulns
}
