package utils

// Gets the first index of a string from an array of strings
func IndexOfStr(list []string, index string) int {
	for i, b := range list {
		if b == index {
			return i
		}
	}
	return -1
}

// Checks to see if a string from an array of strings exists
func ElExistsStr(list []string, index string) bool {
	for _, b := range list {
		if b == index {
			return true
		}
	}
	return false
}
