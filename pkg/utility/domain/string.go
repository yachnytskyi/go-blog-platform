package domain

import (
	"strings"
)

// SanitizeString trims leading and trailing white spaces from the input string.
func SanitizeString(data string) string {
	return strings.TrimSpace(data)
}

// ToLowerString maps the input string to lowercase.
func ToLowerString(data string) string {
	return strings.ToLower(data)
}

// SanitizeAndToLowerString trims leading and trailing white spaces from the input string
func SanitizeAndToLowerString(data string) string {
	return ToLowerString(SanitizeString(data))
}
