package domain

import (
	"fmt"
	"regexp"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

// ValidateField validates a field based on the provided commonValidator.
func ValidateField(logger interfaces.Logger, location, field string, commonValidator domainModel.CommonValidator, validationErrors []error) []error {
	errors := validationErrors

	// Skip validation if the field is optional and empty.
	if commonValidator.IsOptional && field == "" {
		return errors
	}

	// Validate field length.
	if IsStringLengthInvalid(field, commonValidator.MinLength, commonValidator.MaxLength) {
		notification := fmt.Sprintf(constants.StringAllowedLength, commonValidator.MinLength, commonValidator.MaxLength)
		if commonValidator.IsOptional {
			notification = fmt.Sprintf(constants.StringOptionalAllowedLength, commonValidator.MinLength, commonValidator.MaxLength)
		}

		commonValidator.Notification = notification
		validationError := domainError.NewValidationError(
			location+".ValidateField.IsStringLengthInvalid",
			commonValidator.FieldName,
			fieldRequirement(commonValidator.IsOptional),
			commonValidator.Notification,
		)
		logger.Info(validationError)
		errors = append(errors, validationError)
	}

	// Validate field characters.
	if AreStringCharactersInvalid(field, commonValidator.FieldRegex) {
		validationError := domainError.NewValidationError(
			location+".ValidateField.AreStringCharactersInvalid",
			commonValidator.FieldName,
			fieldRequirement(commonValidator.IsOptional),
			commonValidator.Notification,
		)
		logger.Info(validationError)
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
func AreStringCharactersInvalid(checkedString, regexString string) bool {
	if regexp.MustCompile(regexString).MatchString(checkedString) {
		return false
	}

	return true
}
