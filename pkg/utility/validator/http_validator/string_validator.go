package http_validator

import (
	"regexp"
	"strings"
)

func SanitizeString(preparedString *string) {
	*preparedString = strings.TrimSpace(*preparedString)
}

func IsStringNull(checkedString string) bool {
	flag := false

	if checkedString == "" {
		flag = true
	}

	return flag
}

func IsStringLengthExceeded(checkedString string) bool {
	flag := false

	if len(checkedString) > 100 {
		flag = true
	}

	return flag
}

func IsTextStringLengthExceeded(checkedString string) bool {
	flag := false

	if len(checkedString) > 10000 {
		flag = true
	}

	return flag
}

func IsLongStringLengthExceeded(checkedString string) bool {
	flag := false

	if len(checkedString) > 20000 {
		flag = true
	}

	return flag
}

func IsStringContainsSpecialCharacters(checkedString string) bool {
	flag := false

	if !regexp.MustCompile(`^[a-zA-z0-9 !@#$€%^&*{}|()=/\;:+-_~'"<>,.? \t]*$`).MatchString(checkedString) {
		flag = true
	}

	return flag
}

func IsTitleStringContainsSpecialCharacters(checkedString string) bool {
	flag := false

	if !regexp.MustCompile(`^[a-zA-z0-9 !()=[]:;+-_~'",.? \t]*$`).MatchString(checkedString) {
		flag = true
	}

	return flag
}

func IsTextStringContainsSpecialCharacters(checkedString string) bool {
	flag := false

	if !regexp.MustCompile(`^[a-zA-z0-9 !@#$€%^&*{}][|/\()=/\;:+-_~'"<>,.? \t]*$`).MatchString(checkedString) {
		flag = true
	}

	return flag
}
