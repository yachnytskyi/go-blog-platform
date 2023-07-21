package domain_validator

import (
	"regexp"
	"strings"
)

func SanitizeString(preparedString *string) {
	*preparedString = strings.TrimSpace(*preparedString)
}

func IsCorrectLengthText(text string, minLength int, maxLength int) bool {
	flag := true

	if len(text) < minLength || len(text) > maxLength {
		flag = false
	}

	return flag
}

func IsStringContainsSpecialCharacters(checkedString string, regexString string) bool {
	flag := false

	if !regexp.MustCompile(regexString).MatchString(checkedString) {
		flag = true
	}

	return flag
}

func StringsMatch(firstString string, secondString string) bool {
	flag := true

	if firstString != secondString {
		flag = false
	}

	return flag
}

func CheckOptionalStringIsCorrectLengthText(text string, minLength int, maxLength int) bool {
	flag := true

	if len(text) < minLength || len(text) > maxLength {
		flag = false
	}

	return flag
}

func CheckOptionalStringIsStringContainsSpecialCharacters(checkedString string, regexString string) bool {
	flag := false

	if !regexp.MustCompile(regexString).MatchString(checkedString) {
		flag = true
	}

	return flag
}
