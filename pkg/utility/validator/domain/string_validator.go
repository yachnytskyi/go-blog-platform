package domain

import (
	"fmt"
	"regexp"
	"strings"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
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

func ValidateField(field string, commonValidator CommonValidator, validationErrors []error) []error {
	// Initialize a slice to hold validation errors.
	errors := validationErrors
	if IsStringLengthInvalid(field, commonValidator.MinLength, commonValidator.MaxLength) {
		commonValidator.Notification = fmt.Sprintf(constants.StringAllowedLength, commonValidator.MinLength, commonValidator.MaxLength)
		validationError := domainError.NewValidationError(location+"ValidateField.IsStringLengthInvalid",
			commonValidator.FieldName, constants.FieldRequired, commonValidator.Notification)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}
	if AreStringCharactersInvalid(field, commonValidator.FieldRegex) {
		validationError := domainError.NewValidationError(location+"ValidateField.AreStringCharactersInvalid",
			commonValidator.FieldName, constants.FieldRequired, commonValidator.Notification)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	return errors
}

func ValidateOptionalField(field string, commonValidator CommonValidator, validationErrors []error) []error {
	errors := make([]error, len(validationErrors))
	errors = append(errors, validationErrors...)

	if IsStringLengthInvalid(field, commonValidator.MinLength, commonValidator.MaxLength) {
		commonValidator.Notification = fmt.Sprintf(constants.StringOptionalAllowedLength, constants.MaxStringLength)
		validationError := domainError.NewValidationError(location+"ValidateOptionalField.IsStringLengthInvalid",
			commonValidator.FieldName, constants.FieldOptional, commonValidator.Notification)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	if AreStringCharactersInvalid(field, commonValidator.FieldRegex) {
		validationError := domainError.NewValidationError(location+"ValidateField.AreStringCharactersInvalid",
			commonValidator.FieldName, constants.FieldOptional, commonValidator.Notification)
		logging.Logger(validationError)
		errors = append(errors, validationError)
		return errors
	}

	return errors
}

func IsStringLengthInvalid(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}

	return false
}

func AreStringCharactersInvalid(checkedString string, regexString string) bool {
	return !regexp.MustCompile(regexString).MatchString(checkedString)
}

func IsOptionalStringLengthInvalid(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}

	return false
}
