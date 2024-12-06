package domain

import (
	"fmt"
	"regexp"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

// ValidateField validates a field based on the provided StringValidator.
func ValidateField(logger interfaces.Logger, location string, stringValidator StringValidator, validationErrors []error) []error {
	var notification string
	errors := validationErrors

	// Skip validation if the field is optional and empty.
	if stringValidator.IsOptional && stringValidator.Field == "" {
		return errors
	}

	// Validate field length.
	if IsStringLengthInvalid(stringValidator.Field, stringValidator.MinLength, stringValidator.MaxLength) {
		if stringValidator.IsOptional {
			notification = fmt.Sprintf(constants.StringOptionalAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
		} else {
			notification = fmt.Sprintf(constants.StringAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
		}
		validationError := domain.NewValidationError(
			location+".ValidateField.IsStringLengthInvalid",
			stringValidator.FieldName,
			fieldRequirement(stringValidator.IsOptional),
			notification,
		)
		logger.Debug(validationError)
		errors = append(errors, validationError)
	}

	// Validate field characters.
	if AreStringCharactersInvalid(stringValidator.Field, stringValidator.FieldRegex) {
		if stringValidator.Notification == "" {
			stringValidator.Notification = constants.StringAllowedCharacters

		}
		validationError := domain.NewValidationError(
			location+".ValidateField.AreStringCharactersInvalid",
			stringValidator.FieldName,
			fieldRequirement(stringValidator.IsOptional),
			stringValidator.Notification,
		)
		logger.Debug(validationError)
		errors = append(errors, validationError)
	}

	return errors
}

// fieldRequirement returns the field requirement string based on whether the field is optional.
func fieldRequirement(isOptional bool) string {
	if isOptional {
		return constants.FieldOptional
	}

	return constants.FieldRequired
}

// IsStringLengthInvalid checks if the length of the input string is outside the specified range.
func IsStringLengthInvalid(checkedString string, minLength int, maxLength int) bool {
	if len(checkedString) < minLength || len(checkedString) > maxLength {
		return true
	}

	return false
}

// AreStringCharactersInvalid checks if the characters in the input string match the specified regex pattern.
func AreStringCharactersInvalid(checkedString string, regex *regexp.Regexp) bool {
	if regex.MatchString(checkedString) {
		return false
	}

	return true
}
