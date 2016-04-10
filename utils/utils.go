package utils

func IndexOfStr(list []string, index string) int {
	for i, b := range list {
		if b == index {
			return i
		}
	}
	return -1
}

func ElExistsStr(list []string, index string) bool {
	for _, b := range list {
		if b == index {
			return true
		}
	}
	return false
}
