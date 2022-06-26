package helpers

// InArrayString Checks if a value exists in an array
func InArrayString(s string, haystack []string) bool {
	for _, str := range haystack {
		if s == str {
			return true
		}
	}
	return false
}

// ArrayInArrayString Checks if an array exists in array
func ArrayInArrayString(find []string, haystack []string) bool {
	for _, str := range haystack {
		for _, s := range find {
			if s == str {
				return true
			}
		}
	}
	return false
}
