package utility

import "strings"

const (
	firstElement = 0
)

func UserFirstName(firstName string) string {
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[firstElement]
	}

	return firstName
}
