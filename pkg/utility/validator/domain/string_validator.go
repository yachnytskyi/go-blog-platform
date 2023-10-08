package domain

import (
	"fmt"
	"regexp"
	"strings"

	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
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

func ValidateField(field, fieldName, fieldRegex, notification string) domainError.ValidationError {
	if IsStringLengthNotValid(field, constant.MinStringLength, constant.MaxStringLength) {
		notification = fmt.Sprintf(constant.StringAllowedLength, constant.MinStringLength, constant.MaxStringLength)
		return domainError.NewValidationError(fieldName, constant.FieldRequired, notification)
	}
	if IsStringCharactersNotValid(field, fieldRegex) {
		return domainError.NewValidationError(fieldName, constant.FieldRequired, notification)
	}
	return domainError.ValidationError{}
}

func ValidateOptionalField(field, fieldName, fieldType, fieldRegex, notification string) domainError.ValidationError {
	if IsStringLengthNotValid(field, constant.MinOptionalStringLength, constant.MaxStringLength) {
		notification = fmt.Sprintf(constant.StringOptionalAllowedLength, constant.MaxStringLength)
		return domainError.NewValidationError(fieldName, fieldType, notification)
	}
	if IsStringCharactersNotValid(field, fieldRegex) {
		return domainError.NewValidationError(fieldName, fieldType, notification)
	}
	return domainError.ValidationError{}
}

func IsStringLengthValid(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return false
	}
	return true
}

func IsStringLengthNotValid(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}
	return false
}

func IsStringCharactersValid(checkedString string, regexString string) bool {
	return regexp.MustCompile(regexString).MatchString(checkedString)
}

func IsStringCharactersNotValid(checkedString string, regexString string) bool {
	return !regexp.MustCompile(regexString).MatchString(checkedString)
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
