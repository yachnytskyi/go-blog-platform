package domain

import (
	"fmt"
	"regexp"
	"strings"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	// Error Messages.
	stringAllowedLength = "can be between %d and %d characters long"
)

func SanitizeString(data string) string {
	return strings.TrimSpace(data)
}

func ToLowerString(data string) string {
	return strings.ToLower(data)
}

func SanitizeAndToLowerString(data string) string {
	data = strings.TrimSpace(data)
	return strings.ToLower(data)
}

func ValidateField(field, fieldName, fieldType, fieldRegex, errorMessage string) domainError.ValidationError {
	if IsStringLengthNotCorrect(field, config.MinLength, config.MaxLength) {
		return domainError.NewValidationError(fieldName, fieldType, fmt.Sprintf(stringAllowedLength, config.MinLength, config.MaxLength))
	} else if CheckSpecialCharactersString(field, fieldRegex) {
		return domainError.NewValidationError(fieldName, fieldType, errorMessage)
	}
	return domainError.ValidationError{}
}

func IsStringLengthCorrect(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return false
	}
	return true
}

func IsStringLengthNotCorrect(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}
	return false
}

func CheckSpecialCharactersString(checkedString string, regexString string) bool {
	if validator.IsBooleanNotTrue(regexp.MustCompile(regexString).MatchString(checkedString)) {
		return true
	}
	return false
}

func CheckNoSpecialCharactersString(checkedString string, regexString string) bool {
	if validator.IsBooleanNotTrue(regexp.MustCompile(regexString).MatchString(checkedString)) {
		return false
	}
	return true
}

func CheckCorrectLengthOptionalString(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return false
	}
	return true
}

func CheckIncorrectLengthOptionalString(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}
	return false
}

func CheckSpecialCharactersOptionalString(checkedString string, regexString string) bool {
	if validator.IsBooleanNotTrue(regexp.MustCompile(regexString).MatchString(checkedString)) {
		return true
	}
	return false
}

func CheckNoSpecialCharactersOptionalString(checkedString string, regexString string) bool {
	if validator.IsBooleanNotTrue(regexp.MustCompile(regexString).MatchString(checkedString)) {
		return false
	}
	return true
}
