package domain

import "strings"

// SanitizeString trims leading and trailing white spaces from the input string.
// Parameters:
// - data: The string to be sanitized.
// Returns:
// - The sanitized string.
func SanitizeString(data string) string {
	return strings.TrimSpace(data)
}

// ToLowerString converts the input string to lowercase.
// Parameters:
// - data: The string to be converted to lowercase.
// Returns:
// - The lowercase string.
func ToLowerString(data string) string {
	return strings.ToLower(data)
}

// SanitizeAndToLowerString trims leading and trailing white spaces from the input string
// and converts it to lowercase.
// Parameters:
// - data: The string to be sanitized and converted to lowercase.
// Returns:
// - The sanitized and lowercase string.
func SanitizeAndToLowerString(data string) string {
	return strings.ToLower(strings.TrimSpace(data))
}
