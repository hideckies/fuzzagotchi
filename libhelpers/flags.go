package libhelpers

import (
	"strings"
)

type Flags struct {
	Color           bool
	ContentLength   string
	NoContentLength string
	Cookie          string
	Header          string
	Method          string
	PostData        string
	Rate            string
	Status          []int
	HideStatus      []int
	Threads         int
	Timeout         int
	Url             string
	Verbose         bool
	Wordlist        string
}

func (f *Flags) ValidateEGG() bool {
	if strings.Contains(f.Cookie, "EGG") {
		return true
	}
	if strings.Contains(f.Header, "EGG") {
		return true
	}
	if strings.Contains(f.Method, "EGG") {
		return true
	}
	if strings.Contains(f.PostData, "EGG") {
		return true
	}
	if strings.Contains(f.Url, "EGG") {
		return true
	}
	return false
}

func NewFlags() Flags {
	var flags Flags
	flags.Color = true
	flags.ContentLength = "-1"
	flags.NoContentLength = "-1"
	flags.Cookie = ""
	flags.Header = ""
	flags.Method = ""
	flags.PostData = ""
	flags.Rate = ""
	flags.Status = []int{}
	flags.HideStatus = []int{}
	flags.Threads = 0
	flags.Timeout = 0
	flags.Url = ""
	flags.Verbose = false
	flags.Wordlist = ""

	return flags
}
