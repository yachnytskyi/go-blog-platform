package domain

import (
	"fmt"
	"regexp"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location = "pkg.utility.validator.domain."
)

// ValidateField validates a required field based on the provided commonValidator.
func ValidateField(logger applicationModel.Logger, location, field string, commonValidator domainModel.CommonValidator, validationErrors []error) []error {
	errors := validationErrors

	if IsStringLengthInvalid(field, commonValidator.MinLength, commonValidator.MaxLength) {
		commonValidator.Notification = fmt.Sprintf(constants.StringAllowedLength, commonValidator.MinLength, commonValidator.MaxLength)
		validationError := domainError.NewValidationError(
			location+".ValidateField.IsStringLengthInvalid",
			commonValidator.FieldName,
			constants.FieldRequired,
			commonValidator.Notification,
		)

		logger.Warn(validationError)
		errors = append(errors, validationError)
		return errors
	}

	// Check if the string characters are invalid based on the regex pattern.
	if AreStringCharactersInvalid(field, commonValidator.FieldRegex) {
		validationError := domainError.NewValidationError(
			location+".ValidateField.AreStringCharactersInvalid",
			commonValidator.FieldName,
			constants.FieldRequired,
			commonValidator.Notification,
		)

		logger.Warn(validationError)
		errors = append(errors, validationError)
		return errors
	}

	return errors
}

// ValidateOptionalField validates an optional field based on the provided commonValidator.
func ValidateOptionalField(logger applicationModel.Logger, location, field string, commonValidator domainModel.CommonValidator, validationErrors []error) []error {
	if len(field) == 0 {
		return validationErrors
	}

	errors := validationErrors
	if IsStringLengthInvalid(field, commonValidator.MinLength, commonValidator.MaxLength) {
		commonValidator.Notification = fmt.Sprintf(constants.StringOptionalAllowedLength, commonValidator.MinLength, commonValidator.MaxLength)
		validationError := domainError.NewValidationError(
			location+".ValidateOptionalField.IsStringLengthInvalid",
			commonValidator.FieldName,
			constants.FieldOptional,
			commonValidator.Notification,
		)

		logger.Warn(validationError)
		errors = append(errors, validationError)
		return errors
	}

	// Check if the string characters are invalid based on the regex pattern.
	if AreStringCharactersInvalid(field, commonValidator.FieldRegex) {
		validationError := domainError.NewValidationError(
			location+".ValidateOptionalField.AreStringCharactersInvalid",
			commonValidator.FieldName,
			constants.FieldOptional,
			commonValidator.Notification,
		)

		logger.Warn(validationError)
		errors = append(errors, validationError)
		return errors
	}

	return errors
}

// IsStringLengthInvalid checks if the length of the input string is outside the specified range.
func IsStringLengthInvalid(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}

	return false
}

// AreStringCharactersInvalid checks if the characters in the input string match the specified regex pattern.
func AreStringCharactersInvalid(checkedString, regexString string) bool {
	if regexp.MustCompile(regexString).MatchString(checkedString) {
		return false
	}

	return true
}
