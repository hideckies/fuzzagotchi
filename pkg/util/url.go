package util

import "strings"

// Adjust the suffix for the given URL
func AdjustUrlSuffix(url string) string {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url
}
