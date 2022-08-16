package libgotchi

import "strings"

func AdjustUrlSuffix(url string) string {
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return url
}
