package libutils

func IntContains(n []int, num int) bool {
	for _, v := range n {
		if v == num {
			return true
		}
	}

	return false
}

func StringContains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}

	return false
}
