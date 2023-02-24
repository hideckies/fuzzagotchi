package fuzzer

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/hideckies/fuzzagotchi/pkg/output"
)

// Get accurate delay from flag value
func getDelay(rateStr string) time.Duration {
	var delay time.Duration

	r, _ := regexp.Compile("[+]?([0-9]*[.])?[0-9]+")
	rrange, _ := regexp.Compile("([+]?([0-9]*[.])?[0-9]+)-([+]?([0-9]*[.])?[0-9]+)")

	if rrange.MatchString(rateStr) {
		durations := strings.Split(rateStr, "-")
		dmin, _ := strconv.ParseFloat(durations[0], 64)
		dmax, _ := strconv.ParseFloat(durations[1], 64)
		if dmin > dmax {
			color.Red(output.ERROR_DELAY)
			os.Exit(0)
		} else if dmin == dmax {
			s := fmt.Sprintf("%fs", dmin)
			delay, _ = time.ParseDuration(s)
		} else if dmin < dmax {
			drand := dmin + rand.Float64()*(dmax-dmin)
			s := fmt.Sprintf("%fs", drand)
			delay, _ = time.ParseDuration(s)
		} else {
			color.Red(output.ERROR_DELAY)
			os.Exit(0)
		}
	} else if r.MatchString(rateStr) {
		s := fmt.Sprintf("%vs", rateStr)
		delay, _ = time.ParseDuration(s)
	} else {
		color.Red(output.ERROR_DELAY)
		os.Exit(0)
	}

	return delay
}

// func ValidateFlagRate(flags helper.Flags) bool {
// 	rate := flags.Rate
// 	_ = rate
// 	return true
// }
