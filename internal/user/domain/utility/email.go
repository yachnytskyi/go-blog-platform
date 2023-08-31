package utility

import "strings"

func UserFirstName(firstName string) string {
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[0]
	}
	return firstName
}
