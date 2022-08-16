package libgotchi

import (
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hideckies/fuzzagotchi/libhelpers"
)

func NewDuration(flags libhelpers.Flags) time.Duration {
	var duration time.Duration
	if strings.Contains(flags.TimeDelay, "-") {
		durations := strings.Split(flags.TimeDelay, "-")
		dmin, _ := strconv.Atoi(durations[0])
		dmax, _ := strconv.Atoi(durations[1])
		duration = time.Duration(rand.Intn(dmax-dmin)) * time.Millisecond
	} else {
		d, _ := strconv.Atoi(flags.TimeDelay)
		duration = time.Duration(d) * time.Millisecond
	}
	return duration
}

func ValidateFlagTimeDelay(flags libhelpers.Flags) bool {
	timeDelay := flags.TimeDelay
	_ = timeDelay
	return true
}
