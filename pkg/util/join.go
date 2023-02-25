package util

import (
	"strconv"
	"strings"
)

// Join an int array
func IntJoin(intarr []int, sep string) string {
	strarr := make([]string, 0)
	for _, i := range intarr {
		strarr = append(strarr, strconv.Itoa(i))
	}
	return strings.Join(strarr, ",")
}
