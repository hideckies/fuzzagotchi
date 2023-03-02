package fuzzer

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/hideckies/fuzzagotchi/pkg/util"
)

// Check if the result's Content-Length mathces the option's Content-Length
func (f *Fuzzer) matchContentLength(cl int) bool {
	if f.Config.MatchLength == "" && f.Config.HideLength == "" {
		return true
	}

	match := true

	rcl, _ := regexp.Compile("^([1-9][0-9]*|0)$")
	rclrange, _ := regexp.Compile("^(([1-9][0-9]*|0)-([1-9][0-9]*|0))$")

	if f.Config.MatchLength != "" {
		// The config's Content-Length is number only e.g. `120`, `327`, etc.
		if rcl.MatchString(f.Config.MatchLength) {
			mcl, _ := strconv.Atoi(f.Config.MatchLength)
			if cl == mcl {
				return true
			} else {
				return false
			}
		}

		// The config's Content-Length is multiple numbers e.g. '120,123', `320,563,1021`, etc.
		if strings.Contains(f.Config.MatchLength, ",") {
			mcls := strings.Split(f.Config.MatchLength, ",")
			if util.ContainString(mcls, strconv.Itoa(cl)) {
				return true
			} else {
				return false
			}
		}

		// The config's Content-Length is number range e.g. `1-200`, `320-560`, etc.
		if rclrange.MatchString(f.Config.MatchLength) {
			var mclMin int
			var mclMax int
			var err error
			mcls := strings.Split(f.Config.MatchLength, "-")
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
	}

	if f.Config.HideLength != "" {
		if rcl.MatchString(f.Config.HideLength) {
			mhcl, _ := strconv.Atoi(f.Config.HideLength)
			if cl == mhcl {
				return false
			} else {
				return true
			}
		}

		if strings.Contains(f.Config.HideLength, ",") {
			mhcls := strings.Split(f.Config.HideLength, ",")
			if util.ContainString(mhcls, strconv.Itoa(cl)) {
				return false
			} else {
				return true
			}
		}

		if rclrange.MatchString(f.Config.HideLength) {
			var mhclMin int
			var mhclMax int
			var err error
			mhcls := strings.Split(f.Config.HideLength, "-")
			mhclMin, err = strconv.Atoi(mhcls[0])
			if err != nil {
				return false
			}
			if len(mhcls) > 1 {
				mhclMax, err = strconv.Atoi(mhcls[1])
				if err != nil {
					return false
				}
			}

			if mhclMin <= cl && cl <= mhclMax {
				return false
			} else {
				return true
			}
		}
	}

	return match
}
