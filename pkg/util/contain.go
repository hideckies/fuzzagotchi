package util

// Check if the int array contains given int
func ContainInt(n []int, num int) bool {
	for _, v := range n {
		if v == num {
			return true
		}
	}
	return false
}

// Check if the string array contains given string
func ContainString(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}
