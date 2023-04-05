package utils

import (
	"encoding/base64"
	"strings"
)

func Encode(baseString string) string {
	data := base64.StdEncoding.EncodeToString([]byte(baseString))
	return data
}

func Decode(encodedString string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(encodedString)

	if err != nil {
		return "", err
	}

	return string(data), nil
}

func UserFirstName(firstName string) string {
	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	return firstName
}
