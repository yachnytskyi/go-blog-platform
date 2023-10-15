package domain

import (
	"fmt"
	"regexp"
	"strings"

	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location = "pkg.utility.validator.domain."
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
		validationError := domainError.NewValidationError(location+"ValidateField.IsStringLengthNotValid", fieldName, constant.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	if IsStringCharactersNotValid(field, fieldRegex) {
		validationError := domainError.NewValidationError(location+"ValidateField.IsStringCharactersNotValid", fieldName, constant.FieldRequired, notification)
		logging.Logger(validationError)
		return validationError
	}
	return domainError.ValidationError{}
}

func ValidateOptionalField(field, fieldName, fieldType, fieldRegex, notification string) domainError.ValidationError {
	if IsStringLengthNotValid(field, constant.MinOptionalStringLength, constant.MaxStringLength) {
		notification = fmt.Sprintf(constant.StringOptionalAllowedLength, constant.MaxStringLength)
		validationError := domainError.NewValidationError(location+"ValidateOptionalField.IsStringLengthNotValid", fieldName, fieldType, notification)
		logging.Logger(validationError)
		return validationError
	}
	if IsStringCharactersNotValid(field, fieldRegex) {
		validationError := domainError.NewValidationError(location+"ValidateField.IsStringCharactersNotValid", fieldName, fieldType, notification)
		logging.Logger(validationError)
		return validationError
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
