package domain

import (
	"regexp"
	"strings"
)

func SanitizeString(preparedString *string) {
	*preparedString = strings.TrimSpace(*preparedString)
}

func StringToLower(preparedString *string) {
	*preparedString = strings.ToLower(*preparedString)
}

func CheckCorrectLengthString(text string, minLength int, maxLength int) bool {
	if len(text) < minLength || len(text) > maxLength {
		return false
	}
	return true
}

func CheckSpecialCharactersString(checkedString string, regexString string) bool {
	if !regexp.MustCompile(regexString).MatchString(checkedString) {
		return true
	}
	return false
}

func CheckCorrectLengthOptionalString(text string, minLength int, maxLength int) bool {
	if len(text) < minLength || len(text) > maxLength {
		return false
	}
	return true
}

func CheckSpecialCharactersOptionalString(checkedString string, regexString string) bool {
	if !regexp.MustCompile(regexString).MatchString(checkedString) {
		return true
	}
	return false
}

func CheckIncorrectLengthString(text string, minLength int, maxLength int) bool {
	if len(text) < minLength || len(text) > maxLength {
		return true
	}
	return false
}

func CheckNoSpecialCharactersString(checkedString string, regexString string) bool {
	if !regexp.MustCompile(regexString).MatchString(checkedString) {
		return false
	}
	return true
}

func CheckNoMatchStrings(firstString string, secondString string) bool {
	if firstString != secondString {
		return true
	}
	return false
}

func CheckIncorrectLengthOptionalString(text string, minLength int, maxLength int) bool {
	if len(text) < minLength || len(text) > maxLength {
		return true
	}
	return false
}

func CheckNoSpecialCharactersOptionalString(checkedString string, regexString string) bool {
	if !regexp.MustCompile(regexString).MatchString(checkedString) {
		return false
	}
	return true
}
