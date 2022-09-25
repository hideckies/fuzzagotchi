package libgotchi

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hideckies/fuzzagotchi/libhelpers"
)

func NewRate(rateStr string) time.Duration {
	var rate time.Duration

	r, _ := regexp.Compile("[+]?([0-9]*[.])?[0-9]+")
	rrange, _ := regexp.Compile("([+]?([0-9]*[.])?[0-9]+)-([+]?([0-9]*[.])?[0-9]+)")

	if rrange.MatchString(rateStr) {
		durations := strings.Split(rateStr, "-")
		dmin, _ := strconv.ParseFloat(durations[0], 64)
		dmax, _ := strconv.ParseFloat(durations[1], 64)
		if dmin > dmax {
			fmt.Println(ERROR_RATE)
			os.Exit(0)
		} else if dmin == dmax {
			s := fmt.Sprintf("%fs", dmin)
			rate, _ = time.ParseDuration(s)
		} else if dmin < dmax {
			drand := dmin + rand.Float64()*(dmax-dmin)
			s := fmt.Sprintf("%fs", drand)
			rate, _ = time.ParseDuration(s)
		} else {
			fmt.Println(ERROR_RATE)
			os.Exit(0)
		}
	} else if r.MatchString(rateStr) {
		s := fmt.Sprintf("%vs", rateStr)
		rate, _ = time.ParseDuration(s)
	} else {
		fmt.Println(ERROR_RATE)
		os.Exit(0)
	}

	return rate
}

func ValidateFlagRate(flags libhelpers.Flags) bool {
	rate := flags.Rate
	_ = rate
	return true
}
