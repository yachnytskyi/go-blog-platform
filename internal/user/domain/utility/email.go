package utility

import "strings"

// UserFirstName extracts the first name from a given string.
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
