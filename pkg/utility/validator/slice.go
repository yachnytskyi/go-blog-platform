package validator

// IsSliceNotEmpty checks if a slice is not empty.
// Parameters:
// - s: The slice to check.
// Returns:
// - A boolean indicating whether the slice is not empty.
func IsSliceNotEmpty[T any](s []T) bool {
	return len(s) > 0
}

// IsSliceContains checks if a value is present in a slice of strings.
// Parameters:
// - slice: The slice of strings to check.
// - value: The string value to check for existence in the slice.
// Returns:
// - A boolean indicating whether the value is present in the slice.
func IsSliceContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return true
		}
	}
	return false
}

// IsSliceNotContains checks if a value is not present in a slice of strings.
// Parameters:
// - slice: The slice of strings to check.
// - value: The string value to check for non-existence in the slice.
// Returns:
// - A boolean indicating whether the value is not present in the slice.
func IsSliceNotContains(slice []string, value string) bool {
	for _, item := range slice {
		if item == value {
			return false
		}
	}

	return true
}
