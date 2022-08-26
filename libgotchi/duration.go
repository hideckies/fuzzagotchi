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

func NewDuration(flags libhelpers.Flags) time.Duration {
	var duration time.Duration

	r, _ := regexp.Compile("[+]?([0-9]*[.])?[0-9]+")
	rrange, _ := regexp.Compile("([+]?([0-9]*[.])?[0-9]+)-([+]?([0-9]*[.])?[0-9]+)")

	if rrange.MatchString(flags.TimeDelay) {
		durations := strings.Split(flags.TimeDelay, "-")
		dmin, _ := strconv.ParseFloat(durations[0], 64)
		dmax, _ := strconv.ParseFloat(durations[1], 64)
		if dmin > dmax {
			fmt.Println(ERROR_DURATION)
			os.Exit(0)
		} else if dmin == dmax {
			s := fmt.Sprintf("%fs", dmin)
			duration, _ = time.ParseDuration(s)
		} else if dmin < dmax {
			drand := dmin + rand.Float64()*(dmax-dmin)
			s := fmt.Sprintf("%fs", drand)
			duration, _ = time.ParseDuration(s)
		} else {
			fmt.Println(ERROR_DURATION)
			os.Exit(0)
		}
	} else if r.MatchString(flags.TimeDelay) {
		s := fmt.Sprintf("%vs", flags.TimeDelay)
		duration, _ = time.ParseDuration(s)
	} else {
		fmt.Println(ERROR_DURATION)
		os.Exit(0)
	}

	return duration
}

func ValidateFlagTimeDelay(flags libhelpers.Flags) bool {
	timeDelay := flags.TimeDelay
	_ = timeDelay
	return true
}
