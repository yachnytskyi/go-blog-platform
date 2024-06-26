package domain

import (
	"fmt"
	"regexp"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location = "pkg.utility.validator.domain."
)

// ValidateField validates a required field based on the provided commonValidator.
// It checks the string length and character validity using regular expressions.
// Parameters:
// - field: The field to be validated.
// - commonValidator: The validator containing validation rules.
// - validationErrors: A slice of existing validation errors to be appended to.
// Returns:
// - A slice of validation errors including any new errors found.
func ValidateField(field string, commonValidator domainModel.CommonValidator, validationErrors []error) []error {
	// Initialize a slice to hold validation errors.
	errors := validationErrors

	// Check if the string length is invalid.
	if IsStringLengthInvalid(field, commonValidator.MinLength, commonValidator.MaxLength) {
		// Set the notification message for string length violation.
		commonValidator.Notification = fmt.Sprintf(constants.StringAllowedLength, commonValidator.MinLength, commonValidator.MaxLength)
		// Create a new validation error with context and log it.
		validationError := domainError.NewValidationError(location+"ValidateField.IsStringLengthInvalid",
			commonValidator.FieldName, constants.FieldRequired, commonValidator.Notification)
		logging.Logger(validationError)
		// Append the validation error to the errors slice and return.
		errors = append(errors, validationError)
		return errors
	}

	// Check if the string characters are invalid based on the regex pattern.
	if AreStringCharactersInvalid(field, commonValidator.FieldRegex) {
		// Create a new validation error with context and log it.
		validationError := domainError.NewValidationError(location+"ValidateField.AreStringCharactersInvalid",
			commonValidator.FieldName, constants.FieldRequired, commonValidator.Notification)
		logging.Logger(validationError)
		// Append the validation error to the errors slice and return.
		errors = append(errors, validationError)
		return errors
	}

	// Return the accumulated validation errors.
	return errors
}

// ValidateOptionalField validates an optional field based on the provided commonValidator.
// It checks the string length and character validity using regular expressions.
// Parameters:
// - field: The field to be validated.
// - commonValidator: The validator containing validation rules.
// - validationErrors: A slice of existing validation errors to be appended to.
// Returns:
// - A slice of validation errors including any new errors found.
func ValidateOptionalField(field string, commonValidator domainModel.CommonValidator, validationErrors []error) []error {
	// If the field is empty, no validation is needed for optional fields.
	if len(field) == 0 {
		return validationErrors
	}

	// Initialize a slice to hold validation errors.
	errors := validationErrors

	// Check if the string length is invalid.
	if IsStringLengthInvalid(field, commonValidator.MinLength, commonValidator.MaxLength) {
		// Set the notification message for optional string length violation.
		commonValidator.Notification = fmt.Sprintf(constants.StringOptionalAllowedLength, commonValidator.MinLength, commonValidator.MaxLength)

		// Create a new validation error with context and log it.
		validationError := domainError.NewValidationError(location+"ValidateOptionalField.IsStringLengthInvalid",
			commonValidator.FieldName, constants.FieldOptional, commonValidator.Notification)
		logging.Logger(validationError)
		// Append the validation error to the errors slice and return.
		errors = append(errors, validationError)
		return errors
	}

	// Check if the string characters are invalid based on the regex pattern.
	if AreStringCharactersInvalid(field, commonValidator.FieldRegex) {
		// Create a new validation error with context and log it.
		validationError := domainError.NewValidationError(location+"ValidateOptionalField.AreStringCharactersInvalid",
			commonValidator.FieldName, constants.FieldOptional, commonValidator.Notification)
		logging.Logger(validationError)
		// Append the validation error to the errors slice and return.
		errors = append(errors, validationError)
		return errors
	}

	// Return the accumulated validation errors.
	return errors
}

// IsStringLengthInvalid checks if the length of the input string is outside the specified range.
// Parameters:
// - checkedString: The string to check.
// - minLength: The minimum allowed length.
// - maxLength: The maximum allowed length.
// Returns:
// - A boolean indicating whether the string length is invalid.
func IsStringLengthInvalid(checkedString string, minLength int, maxLength int) bool {
	return len(checkedString) < minLength || len(checkedString) > maxLength
}

// AreStringCharactersInvalid checks if the characters in the input string match the specified regex pattern.
// Parameters:
// - checkedString: The string to check.
// - regexString: The regex pattern to match against.
// Returns:
// - A boolean indicating whether the string characters are invalid.
func AreStringCharactersInvalid(checkedString string, regexString string) bool {
	return !regexp.MustCompile(regexString).MatchString(checkedString)
}
