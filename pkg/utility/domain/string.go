package domain

import (
	"strings"
)

// SanitizeString trims leading and trailing white spaces from the input string.
//
// Parameters:
// - data: The string to be sanitized.
//
// Returns:
// - The sanitized string.
func SanitizeString(data string) string {
	return strings.TrimSpace(data)
}

// ToLowerString maps the input string to lowercase.
//
// Parameters:
// - data: The string to be mapped to lowercase.
//
// Returns:
// - The lowercase string.
func ToLowerString(data string) string {
	return strings.ToLower(data)
}

// SanitizeAndToLowerString trims leading and trailing white spaces from the input string
// and maps it to lowercase by utilizing SanitizeString and ToLowerString functions.
//
// Parameters:
// - data: The string to be sanitized and mapped to lowercase.
//
// Returns:
// - The sanitized and lowercase string.
func SanitizeAndToLowerString(data string) string {
	// Utilize existing functions for clarity and to avoid redundancy
	return ToLowerString(SanitizeString(data))
}
