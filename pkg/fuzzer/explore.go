package fuzzer

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/pkg/output"
	"github.com/hideckies/fuzzagotchi/pkg/util"
)

type ExplorePage struct {
	Body []byte `json:"body"`
	Path string `json:"path"`

	Keywords []string `json:"keywords"`
}

type ExploreResult struct {
	Path     string   `json:"path"`
	Keywords []string `json:"keywords"`
}

type Explore struct {
	Pages []ExplorePage

	Results []ExploreResult
}

// Initialize a new Explore
func NewExplore(resps []Response) Explore {
	var e Explore
	e.Pages = make([]ExplorePage, 0)

	pathsToExplore := make([]string, 0)
	for _, resp := range resps {
		if util.ContainString(pathsToExplore, resp.Path) {
			continue
		}
		pathsToExplore = append(pathsToExplore, resp.Path)

		page := ExplorePage{Body: resp.Body, Path: resp.Path}
		e.Pages = append(e.Pages, page)
	}
	return e
}

// Explore each page to find various information
func (e *Explore) explore() {
	for _, page := range e.Pages {
		reader := bytes.NewReader(page.Body)
		scanner := bufio.NewScanner(reader)
		for scanner.Scan() {
			txt := scanner.Text()
			keywords := e.findKeywords(txt)
			if len(keywords) == 0 {
				continue
			}

			result := ExploreResult{Path: page.Path, Keywords: keywords}
			e.Results = append(e.Results, result)
		}
	}

	e.printResult()
}

// Find keywords about sensitive information
func (e *Explore) findKeywords(txt string) []string {
	keywords := make([]string, 0)

	// Credentials
	reCreds := regexp.MustCompile(`/user|User|USER|pass|Pass|PASS/g`)
	mCreds := reCreds.FindAllString(txt, -1)
	if len(mCreds) > 0 {
		keywords = append(keywords, mCreds...)
	}

	// Base64
	reBase := regexp.MustCompile(`/[a-zA-Z0-9]+\=\=/g`)
	mBase := reBase.FindAllString(txt, -1)
	if len(mBase) > 0 {
		keywords = append(keywords, mBase...)
	}

	return keywords
}

// Find vulnerabilities
func (e *Explore) findVuln(txt string) {}

// Print result
func (e *Explore) printResult() {
	// Keywords
	fmt.Printf("%s\n", color.YellowString(output.TMPL_BAR_DOUBLE_M))
	fmt.Printf("%s %s\n", color.CyanString("+"), color.CyanString("KEYWORDS IN PAGES"))
	fmt.Printf("%s\n", color.YellowString(output.TMPL_BAR_DOUBLE_M))
	if len(e.Results) > 0 {
		for _, result := range e.Results {
			if len(result.Keywords) > 0 {
				color.White("%s:\n", result.Path)
				for _, keyword := range result.Keywords {
					color.Cyan(" %s\n", keyword)
				}
			}
		}
	}
}
