package utility

import "strings"

// UserFirstName extracts the first name from a given string.
// If the string contains spaces, it splits and returns the first word.
// Otherwise, it returns the original string.
//
// Parameters:
// - fullName: A string containing the full name from which the first name will be extracted.
//
// Returns:
// - A string containing the first name or the original string if no spaces are found.
func UserFirstName(fullName string) string {
	// Check if there is a space in the name
	spaceIndex := strings.Index(fullName, " ")
	if spaceIndex != -1 {
		// Return the substring up to the first space.
		return fullName[:spaceIndex]
	}

	// If no space is found, return the original string.
	return fullName
}
