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

func TestValidateField(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "Valid123",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateField", stringValidator, validationErrors)

	assert.Len(t, validationErrors, 0, test.ErrorNilMessage)
}

func TestValidateFieldAtMinLength(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "Min",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldAtMinLength", stringValidator, validationErrors)

	assert.Len(t, validationErrors, 0, test.ErrorNilMessage)
}

func TestValidateFieldAtMaxLength(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "MaximumLen",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldAtMaxLength", stringValidator, validationErrors)

	assert.Len(t, validationErrors, 0, test.ErrorNilMessage)
}

func TestValidateOptionalFieldEmptyField(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: true,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldEmptyField", stringValidator, validationErrors)

	assert.Len(t, validationErrors, 0, test.ErrorNilMessage)
}

func TestValidateFieldTooShort(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "No",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldTooShort", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldTooShort.ValidateField.IsStringLengthInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, expectedNotification)

	assert.Len(t, validationErrors, 1, test.EqualMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldTooLong(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "ThisFieldIsWayTooLong",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldTooLong", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldTooLong.ValidateField.IsStringLengthInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldInvalidCharacters(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	fieldName := "testField"
	stringValidator := utility.StringValidator{
		FieldName:  fieldName,
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "Invalid@!",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldInvalidCharacters", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedCharacters)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldEmptyField(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]*$",
		Field:      "",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldEmptyField", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldEmptyField.ValidateField.IsStringLengthInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldValidCharactersButIncorrectLength(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldRegex: "^[a-zA-Z]+$",
		FieldName:  "testField",
		Field:      "Abc",
		MinLength:  5,
		MaxLength:  10,
		IsOptional: false,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldValidCharactersButIncorrectLength", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldValidCharactersButIncorrectLength.ValidateField.IsStringLengthInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldValidLengthButInvalidCharacters(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	stringValidator := utility.StringValidator{
		FieldRegex: "^[a-zA-Z]+$",
		FieldName:  "testField",
		Field:      "12345",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}

	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldValidLengthButInvalidCharacters", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldValidLengthButInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedCharacters)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldTooShort(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "No",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: true,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldTooShort", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldTooShort.ValidateField.IsStringLengthInvalid"
	expectedNotification := fmt.Sprintf(constants.StringOptionalAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, expectedNotification)

	assert.Len(t, validationErrors, 1, test.EqualMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldTooLong(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "ThisFieldIsWayTooLong",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: true,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldTooLong", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldTooLong.ValidateField.IsStringLengthInvalid"
	expectedNotification := fmt.Sprintf(constants.StringOptionalAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldInvalidCharacters(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	fieldName := "testField"
	stringValidator := utility.StringValidator{
		FieldName:  fieldName,
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "Invalid@!",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: true,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldInvalidCharacters", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedCharacters)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldValidCharactersButIncorrectLength(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldRegex: "^[a-zA-Z]+$",
		FieldName:  "testField",
		Field:      "Abc",
		MinLength:  5,
		MaxLength:  10,
		IsOptional: true,
	}
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldValidCharactersButIncorrectLength", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldValidCharactersButIncorrectLength.ValidateField.IsStringLengthInvalid"
	expectedNotification := fmt.Sprintf(constants.StringOptionalAllowedLength, stringValidator.MinLength, stringValidator.MaxLength)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldValidLengthButInvalidCharacters(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	stringValidator := utility.StringValidator{
		FieldRegex: "^[a-zA-Z]+$",
		FieldName:  "testField",
		Field:      "12345",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: true,
	}

	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldValidLengthButInvalidCharacters", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldValidLengthButInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedNotification := fmt.Sprintf(constants.StringAllowedCharacters)
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, expectedNotification)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldMultipleErrors(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.StringValidator{
		FieldName:  "testField",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: true,
	}
	stringValidator1 := utility.StringValidator{
		FieldName:  "field1",
		FieldRegex: "^[a-zA-Z]+$",
		Field:      "Shor",
		MinLength:  5,
		MaxLength:  10,
		IsOptional: false,
	}
	stringValidator2 := utility.StringValidator{
		FieldName:  "field2",
		FieldRegex: "^[a-zA-Z0-9]+$",
		Field:      "Invalid@!",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}
	stringValidator3 := utility.StringValidator{
		FieldName:  "field3",
		FieldRegex: "^[a-zA-Z0-9]*$",
		Field:      "",
		MinLength:  3,
		MaxLength:  10,
		IsOptional: false,
	}

	validationErrors := make([]error, 0, 3)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldMultipleErrors.Field", stringValidator, validationErrors)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldMultipleErrors.Field1", stringValidator1, validationErrors)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldMultipleErrors.Field2", stringValidator2, validationErrors)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldMultipleErrors.Field3", stringValidator3, validationErrors)

	expectedLocation1 := location + "TestValidateFieldMultipleErrors.Field1.ValidateField.IsStringLengthInvalid"
	expectedNotification1 := fmt.Sprintf(constants.StringAllowedLength, stringValidator1.MinLength, stringValidator1.MaxLength)
	expectedError1 := domain.NewValidationError(expectedLocation1, stringValidator1.FieldName, constants.FieldRequired, expectedNotification1)

	expectedLocation2 := location + "TestValidateFieldMultipleErrors.Field2.ValidateField.AreStringCharactersInvalid"
	expectedNotification2 := fmt.Sprintf(constants.StringAllowedCharacters)
	expectedError2 := domain.NewValidationError(expectedLocation2, stringValidator2.FieldName, constants.FieldRequired, expectedNotification2)

	expectedLocation3 := location + "TestValidateFieldMultipleErrors.Field3.ValidateField.IsStringLengthInvalid"
	expectedNotification3 := fmt.Sprintf(constants.StringAllowedLength, stringValidator3.MinLength, stringValidator3.MaxLength)
	expectedError3 := domain.NewValidationError(expectedLocation3, stringValidator3.FieldName, constants.FieldRequired, expectedNotification3)

	assert.Len(t, validationErrors, 3, test.EqualMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError1, validationErrors[0], test.EqualMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[1], test.EqualMessage)
	assert.Equal(t, expectedError2, validationErrors[1], test.EqualMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[2], test.EqualMessage)
	assert.Equal(t, expectedError3, validationErrors[2], test.EqualMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[2], test.EqualMessage)

}
