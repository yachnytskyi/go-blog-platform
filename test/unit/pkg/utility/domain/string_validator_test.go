package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/domain"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

const (
	location = "test.unit.pkg.model.domain.validator."
)

func TestValidateFieldValidField(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	field := "Valid123"
	stringValidator := utility.StringValidator{
		MinLength:  3,
		MaxLength:  10,
		FieldRegex: "^[a-zA-Z0-9]+$",
		FieldName:  "testField",
		IsOptional: false,
	}
	validationErrors := []error{}
	result := utility.ValidateField(mockLogger, location+"TestValidateFieldValidField", field, stringValidator, validationErrors)

	assert.Len(t, result, 0, test.ErrorNilMessage)
}

func TestValidateFieldTooShort(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	field := "No"
	stringValidator := utility.StringValidator{
		MinLength:  3,
		MaxLength:  10,
		FieldRegex: "^[a-zA-Z0-9]+$",
		FieldName:  "testField",
		IsOptional: false,
	}
	validationErrors := []error{}
	result := utility.ValidateField(mockLogger, location+"TestValidateFieldTooShort", field, stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldTooShort.ValidateField.IsStringLengthInvalid"
	expectedNotification := "Can be between 3 and 10 characters long."
	expectedField := "testField"
	expectedErrorMessage := fmt.Sprintf("location: %s notification: %s field: %s type: %s",
		expectedLocation, expectedNotification, expectedField, constants.FieldRequired)

	actualError := result[0].(domain.ValidationError)
	assert.Len(t, result, 1, test.EqualMessage)
	assert.IsType(t, domain.ValidationError{}, result[0], test.EqualMessage)
	assert.Equal(t, expectedErrorMessage, actualError.Error(), test.EqualMessage)
}

func TestValidateFieldTooLong(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	field := "ThisFieldIsWayTooLong"
	stringValidator := utility.StringValidator{
		MinLength:  3,
		MaxLength:  10,
		FieldRegex: "^[a-zA-Z0-9]+$",
		FieldName:  "testField",
		IsOptional: false,
	}

	validationErrors := []error{}
	result := utility.ValidateField(mockLogger, location+"TestValidateFieldTooLong", field, stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldTooLong.ValidateField.IsStringLengthInvalid"
	expectedNotification := "Can be between 3 and 10 characters long."
	expectedField := "testField"
	expectedErrorMessage := fmt.Sprintf("location: %s notification: %s field: %s type: %s",
		expectedLocation, expectedNotification, expectedField, constants.FieldRequired)

	actualError := result[0].(domain.ValidationError)
	assert.Len(t, result, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, result[0], test.EqualMessage)
	assert.Equal(t, expectedErrorMessage, actualError.Error(), test.EqualMessage)
}

func TestValidateFieldInvalidCharacters(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	field := "Invalid@!"
	fieldName := "testField"
	stringValidator := utility.StringValidator{
		MinLength:  3,
		MaxLength:  10,
		FieldRegex: "^[a-zA-Z0-9]+$", // Adjust regex as needed
		FieldName:  fieldName,
		IsOptional: false,
	}

	validationErrors := []error{}
	result := utility.ValidateField(mockLogger, location+"TestValidateFieldInvalidCharacters", field, stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedCharacters)
	expectedField := fieldName
	expectedErrorMessage := fmt.Sprintf("location: %s notification: %s field: %s type: %s",
		expectedLocation, expectedNotification, expectedField, constants.FieldRequired)

	actualError := result[0].(domain.ValidationError)
	assert.Len(t, result, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, result[0], test.EqualMessage)
	assert.Equal(t, expectedErrorMessage, actualError.Error(), test.EqualMessage)
}

// func TestValidateFieldOptionalEmptyField(t *testing.T) {
// 	t.Parallel()
// 	mockLogger := mock.NewMockLogger()
// 	field := ""
// 	stringValidator := utility.StringValidator{
// 		MinLength:  3,
// 		MaxLength:  10,
// 		FieldRegex: "^[a-zA-Z0-9]+$",
// 		FieldName:  "testField",
// 		IsOptional: true,
// 	}

// 	validationErrors := []error{}
// 	result := utility.ValidateField(mockLogger, location+"TestValidateFieldOptionalEmptyField", field, stringValidator, validationErrors)

// 	assert.Len(t, result, 0, test.ErrorNotNilMessage)
// }

// func TestValidateFieldValidCharactersButIncorrectLength(t *testing.T) {
// 	t.Parallel()
// 	mockLogger := mock.NewMockLogger()
// 	field := "Abc"
// 	stringValidator := utility.StringValidator{
// 		MinLength:  5,
// 		MaxLength:  10,
// 		FieldRegex: "^[a-zA-Z]+$",
// 		FieldName:  "testField",
// 		IsOptional: false,
// 	}

// 	validationErrors := []error{}
// 	result := utility.ValidateField(mockLogger, location+"TestValidateFieldValidCharactersButIncorrectLength", field, stringValidator, validationErrors)

// 	expectedLocation := location + "TestValidateFieldValidCharactersButIncorrectLength.ValidateField.IsStringLengthInvalid"
// 	expectedNotification := "Can be between 5 and 10 characters long."
// 	expectedField := "testField"
// 	expectedErrorMessage := fmt.Sprintf("location: %s notification: %s field: %s type: %s",
// 		expectedLocation, expectedNotification, expectedField, constants.FieldRequired)

// 	actualError := result[0].(domain.ValidationError)
// 	assert.Len(t, result, 1, test.ErrorNotNilMessage)
// 	assert.IsType(t, domain.ValidationError{}, result[0], test.EqualMessage)
// 	assert.Equal(t, expectedErrorMessage, actualError.Error(), test.EqualMessage)
// }

// func TestValidateFieldValidLengthButInvalidCharacters(t *testing.T) {
// 	t.Parallel()
// 	mockLogger := mock.NewMockLogger()
// 	field := "12345"
// 	stringValidator := utility.StringValidator{
// 		MinLength:  3,
// 		MaxLength:  10,
// 		FieldRegex: "^[a-zA-Z]+$",
// 		FieldName:  "testField",
// 		IsOptional: false,
// 	}

// 	validationErrors := []error{}
// 	result := utility.ValidateField(mockLogger, location+"TestValidateFieldValidLengthButInvalidCharacters", field, stringValidator, validationErrors)

// 	expectedLocation := location + "TestValidateFieldValidLengthButInvalidCharacters.ValidateField.AreStringCharactersInvalid"
// 	expectedNotification := "field contains invalid characters"
// 	expectedField := "testField"
// 	expectedErrorMessage := fmt.Sprintf("location: %s notification: %s field: %s type: %s",
// 		expectedLocation, expectedNotification, expectedField, constants.FieldRequired)

// 	actualError := result[0].(domain.ValidationError)
// 	assert.Len(t, result, 1, test.ErrorNotNilMessage)
// 	assert.IsType(t, domain.ValidationError{}, result[0], test.EqualMessage)
// 	assert.Equal(t, expectedErrorMessage, actualError.Error(), test.EqualMessage)
// }

// func TestIsStringLengthInvalidValidLength(t *testing.T) {
// 	t.Parallel()
// 	checkedString := "valid"
// 	minLength := 3
// 	maxLength := 10
// 	result := utility.IsStringLengthInvalid(checkedString, minLength, maxLength)

// 	assert.False(t, result, test.NotFailureMessage)
// }

// func TestIsStringLengthInvalidTooShort(t *testing.T) {
// 	t.Parallel()
// 	checkedString := "no"
// 	minLength := 3
// 	maxLength := 10
// 	result := utility.IsStringLengthInvalid(checkedString, minLength, maxLength)

// 	assert.True(t, result, test.FailureMessage)
// }

// func TestIsStringLengthInvalidTooLong(t *testing.T) {
// 	t.Parallel()
// 	checkedString := "thisisaverylongstring"
// 	minLength := 3
// 	maxLength := 10
// 	result := utility.IsStringLengthInvalid(checkedString, minLength, maxLength)

// 	assert.True(t, result, test.FailureMessage)
// }

// func TestIsStringLengthInvalidExactMinLength(t *testing.T) {
// 	t.Parallel()
// 	checkedString := "min"
// 	minLength := 3
// 	maxLength := 10
// 	result := utility.IsStringLengthInvalid(checkedString, minLength, maxLength)

// 	assert.False(t, result, test.NotFailureMessage)
// }

// func TestIsStringLengthInvalidExactMaxLength(t *testing.T) {
// 	t.Parallel()
// 	checkedString := "maximum"
// 	minLength := 3
// 	maxLength := 7
// 	result := utility.IsStringLengthInvalid(checkedString, minLength, maxLength)

// 	assert.False(t, result, test.NotFailureMessage)
// }

// func TestAreStringCharactersInvalidValidCharacters(t *testing.T) {
// 	t.Parallel()
// 	checkedString := "Valid123"
// 	regexPattern := "^[a-zA-Z0-9]+$"
// 	result := utility.AreStringCharactersInvalid(checkedString, regexPattern)

// 	assert.False(t, result, test.NotFailureMessage)
// }

// func TestAreStringCharactersInvalidInvalidCharacters(t *testing.T) {
// 	t.Parallel()
// 	checkedString := "Invalid!@#"
// 	regexPattern := "^[a-zA-Z0-9]+$"

// 	result := utility.AreStringCharactersInvalid(checkedString, regexPattern)

// 	assert.True(t, result, test.FailureMessage)
// }

// func TestAreStringCharactersInvalidEmptyString(t *testing.T) {
// 	t.Parallel()
// 	checkedString := ""
// 	regexPattern := "^[a-zA-Z0-9]+$"
// 	result := utility.AreStringCharactersInvalid(checkedString, regexPattern)

// 	assert.True(t, result, test.FailureMessage)
// }

// func TestAreStringCharactersInvalidComplexRegex(t *testing.T) {
// 	t.Parallel()
// 	checkedString := "abc-123_def"
// 	regexPattern := `^[a-z0-9_\-]+$`
// 	result := utility.AreStringCharactersInvalid(checkedString, regexPattern)

// 	assert.False(t, result, test.NotFailureMessage)
// }
