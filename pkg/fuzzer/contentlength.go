package fuzzer

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/hideckies/fuzzagotchi/pkg/util"
)

// Check if the result's Content-Length mathces the option's Content-Length
func (f *Fuzzer) matchContentLength(cl int) bool {
	if f.Config.ContentLength == "" {
		return true
	}

	match := true

	rcl, _ := regexp.Compile("^([1-9][0-9]*|0)$")
	rclrange, _ := regexp.Compile("^(([1-9][0-9]*|0)-([1-9][0-9]*|0))$")

	// The config's Content-Length is number only e.g. `120`, `327`, etc.
	if rcl.MatchString(f.Config.ContentLength) {
		mcl, _ := strconv.Atoi(f.Config.ContentLength)
		if cl == mcl {
			return true
		} else {
			return false
		}
	}

	// The config's Content-Length is multiple numbers e.g. '120,123', `320,563,1021`, etc.
	if strings.Contains(f.Config.ContentLength, ",") {
		mcls := strings.Split(f.Config.ContentLength, ",")
		if util.ContainString(mcls, strconv.Itoa(cl)) {
			return true
		} else {
			return false
		}
	}

	// The config's Content-Length is number range e.g. `1-200`, `320-560`, etc.
	if rclrange.MatchString(f.Config.ContentLength) {
		var mclMin int
		var mclMax int
		var err error
		mcls := strings.Split(f.Config.ContentLength, "-")
		mclMin, err = strconv.Atoi(mcls[0])
		if err != nil {
			return false
		}
		if len(mcls) > 1 {
			mclMax, err = strconv.Atoi(mcls[1])
			if err != nil {
				return false
			}
		}

		if mclMin <= cl && cl <= mclMax {
			return true
		} else {
			return false
		}
	}

	return match
}
