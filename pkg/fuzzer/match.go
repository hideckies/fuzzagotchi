package fuzzer

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/hideckies/fuzzagotchi/pkg/util"
)

// Check if the result's status code matches the option's status code
func (f *Fuzzer) matchStatusCode(sc int) bool {
	// adjust match&hide status code
	if f.Config.HideStatus != "" {
		f.Config.MatchStatus = ""
	}
	return checkMatched(sc, f.Config.MatchStatus, f.Config.HideStatus)
}

// Check if the result's Content-Length mathces the option's Content-Length
func (f *Fuzzer) matchContentLength(cl int) bool {
	return checkMatched(cl, f.Config.MatchLength, f.Config.HideLength)
}

// Check if the result's Content-Length mathces the option's Content-Length
func (f *Fuzzer) matchContentWords(words int) bool {
	return checkMatched(words, f.Config.MatchWords, f.Config.HideWords)
}

// Utitlity: check if matched
func checkMatched(num int, mflag string, hflag string) bool {
	if mflag == "" && hflag == "" {
		return true
	}

	match := true

	reFlag, _ := regexp.Compile("^([1-9][0-9]*|0)$")
	reFlagRange, _ := regexp.Compile("^(([1-9][0-9]*|0)-([1-9][0-9]*|0))$")

	if mflag != "" {
		// The config's Content-Length is number only e.g. `120`, `327`, etc.
		if reFlag.MatchString(mflag) {
			mint, _ := strconv.Atoi(mflag)
			if num == mint {
				return true
			} else {
				return false
			}
		}

		// The config's Content-Length is multiple numbers e.g. '120,123', `320,563,1021`, etc.
		if strings.Contains(mflag, ",") {
			mints := strings.Split(mflag, ",")
			if util.ContainString(mints, strconv.Itoa(num)) {
				return true
			} else {
				return false
			}
		}

		// The config's Content-Length is number range e.g. `1-200`, `320-560`, etc.
		if reFlagRange.MatchString(mflag) {
			var mintsMin int
			var mintsMax int
			var err error
			mints := strings.Split(mflag, "-")
			mintsMin, err = strconv.Atoi(mints[0])
			if err != nil {
				return false
			}
			if len(mints) > 1 {
				mintsMax, err = strconv.Atoi(mints[1])
				if err != nil {
					return false
				}
			}

			if mintsMin <= num && num <= mintsMax {
				return true
			} else {
				return false
			}
		}
	}

	if hflag != "" {
		if reFlag.MatchString(hflag) {
			hint, _ := strconv.Atoi(hflag)
			if num == hint {
				return false
			} else {
				return true
			}
		}

		if strings.Contains(hflag, ",") {
			hints := strings.Split(hflag, ",")
			if util.ContainString(hints, strconv.Itoa(num)) {
				return false
			} else {
				return true
			}
		}

		if reFlagRange.MatchString(hflag) {
			var hintsMin int
			var hintsMax int
			var err error
			hints := strings.Split(hflag, "-")
			hintsMin, err = strconv.Atoi(hints[0])
			if err != nil {
				return false
			}
			if len(hints) > 1 {
				hintsMax, err = strconv.Atoi(hints[1])
				if err != nil {
					return false
				}
			}

			if hintsMin <= num && num <= hintsMax {
				return false
			} else {
				return true
			}
		}
	}

	return match
}

// Check if the results's content contains given keyword
func (f *Fuzzer) matchKeyword(content string) bool {
	match := true

	if f.Config.MatchKeyword == "" && f.Config.HideKeyword == "" {
		return true
	}
	if f.Config.MatchKeyword != "" {
		if strings.Contains(content, f.Config.MatchKeyword) {
			return true
		} else {
			return false
		}
	}
	if f.Config.HideKeyword != "" {
		if strings.Contains(content, f.Config.HideKeyword) {
			return false
		} else {
			return true
		}
	}
	return match
}
