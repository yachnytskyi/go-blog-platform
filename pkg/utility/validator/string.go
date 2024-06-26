package validator

// IsStringNotEmpty checks if a string is not empty.
//
// Parameters:
// - data: The string to check.
//
// Returns:
// - A boolean indicating whether the string is not empty.
func IsStringNotEmpty(data string) bool {
	return data != ""
}

// AreStringsNotEqual checks if two strings are not equal.
//
// Parameters:
// - firstString: The first string to compare.
// - secondString: The second string to compare.
//
// Returns:
// - A boolean indicating whether the two strings are not equal.
func AreStringsNotEqual(firstString, secondString string) bool {
	return firstString != secondString
}
