package validator

func IsSliceNotEmpty[T any](s []T) bool {
	return len(s) > 0
}

// IsSliceContains checks if a string is present in a slice of strings.
func IsSliceContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// IsSliceNotContains checks if a string is present in a slice of strings.
func IsSliceNotContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return false
		}
	}
	return true
}
