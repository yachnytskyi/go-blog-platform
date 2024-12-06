package domain

import (
	"fmt"
	"regexp"
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

	testField = "testField"
)

var (
	alphaNumericRegex         = regexp.MustCompile(`^[a-zA-Z0-9]+$`)
	alphaRegex                = regexp.MustCompile(`^[a-zA-Z]+$`)
	alphaNumericOptionalRegex = regexp.MustCompile(`^[a-zA-Z0-9]*$`)
)

func TestValidateField(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	
	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "Valid123", alphaNumericRegex, 3, 10, false)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateField", stringValidator, validationErrors)

	assert.Len(t, validationErrors, 0, test.ErrorNilMessage)
}

func TestValidateFieldAtMinLength(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	
	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "Min", alphaNumericRegex, 3, 10, false)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldAtMinLength", stringValidator, validationErrors)
	assert.Len(t, validationErrors, 0, test.ErrorNilMessage)
}

func TestValidateFieldAtMaxLength(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	
	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "MaximumLen", alphaNumericRegex, 3, 10, false)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldAtMaxLength", stringValidator, validationErrors)

	assert.Len(t, validationErrors, 0, test.ErrorNilMessage)
}

func TestValidateOptionalFieldEmptyField(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "No", alphaNumericRegex, 3, 10, false)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldEmptyField", stringValidator, validationErrors)

	assert.Len(t, validationErrors, 0, test.ErrorNilMessage)
}

func TestValidateFieldTooShort(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "No", alphaNumericRegex, 3, 10, false)
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

	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "ThisFieldIsWayTooLong", alphaNumericRegex, 3, 10, false)
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
	
	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "Invalid@!", alphaNumericRegex, 3, 10, false)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldInvalidCharacters", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, constants.StringAllowedCharacters)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldEmptyField(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	
	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "", alphaNumericOptionalRegex, 3, 10, false)
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
	
	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "Abc", alphaRegex, 5, 10, false)
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

	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "12345", alphaRegex, 3, 10, false)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldValidLengthButInvalidCharacters", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldValidLengthButInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, constants.StringAllowedCharacters)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldWithLeadingAndTrailingSpaces(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "  Valid123  ", alphaNumericRegex, 3, 20, false)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldWithLeadingAndTrailingSpaces", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldWithLeadingAndTrailingSpaces.ValidateField.AreStringCharactersInvalid"
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, constants.StringAllowedCharacters)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldWithInternalTab(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	validationErrors := []error{}
	stringValidator := utility.NewStringValidator(testField, "Valid 123", alphaNumericRegex, 3, 20, false)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldWithLeadingAndTrailingSpaces", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateFieldWithLeadingAndTrailingSpaces.ValidateField.AreStringCharactersInvalid"
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldRequired, constants.StringAllowedCharacters)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldTooShort(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.NewStringValidator(testField, "No", alphaNumericRegex, 3, 10, true)
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

	stringValidator := utility.NewStringValidator(testField, "ThisFieldIsWayTooLong", alphaNumericRegex, 3, 10, true)
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

	stringValidator := utility.NewStringValidator(testField, "Invalid@!", alphaNumericRegex, 3, 10, true)
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldInvalidCharacters", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, constants.StringAllowedCharacters)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldValidCharactersButIncorrectLength(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.NewStringValidator(testField, "Abc", alphaRegex, 5, 10, true)
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

	stringValidator := utility.NewStringValidator(testField, "12345", alphaRegex, 3, 10, true)
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldValidLengthButInvalidCharacters", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldValidLengthButInvalidCharacters.ValidateField.AreStringCharactersInvalid"
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, constants.StringAllowedCharacters)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldWithLeadingAndTrailingSpaces(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.NewStringValidator(testField, "  Valid123  ", alphaNumericRegex, 3, 20, true)
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldWithLeadingAndTrailingSpaces", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldWithLeadingAndTrailingSpaces.ValidateField.AreStringCharactersInvalid"
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, constants.StringAllowedCharacters)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateOptionalFieldWithInternalTab(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.NewStringValidator(testField, "Valid 123", alphaNumericRegex, 3, 20, true)
	validationErrors := []error{}
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateOptionalFieldWithInternalTab", stringValidator, validationErrors)

	expectedLocation := location + "TestValidateOptionalFieldWithInternalTab.ValidateField.AreStringCharactersInvalid"
	expectedError := domain.NewValidationError(expectedLocation, stringValidator.FieldName, constants.FieldOptional, constants.StringAllowedCharacters)

	assert.Len(t, validationErrors, 1, test.ErrorNotNilMessage)
	assert.IsType(t, domain.ValidationError{}, validationErrors[0], test.EqualMessage)
	assert.Equal(t, expectedError, validationErrors[0], test.EqualMessage)
}

func TestValidateFieldMultipleErrors(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	stringValidator := utility.NewStringValidator(testField, "", alphaNumericRegex, 3, 10, true)
	stringValidator1 := utility.NewStringValidator(testField+"1", "Shor", alphaRegex, 5, 10, false)
	stringValidator2 := utility.NewStringValidator(testField+"2", "Invalid@!", alphaNumericRegex, 3, 10, false)
	stringValidator3 := utility.NewStringValidator(testField+"3", "", alphaNumericOptionalRegex, 3, 10, false)

	validationErrors := make([]error, 0, 3)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldMultipleErrors.Field", stringValidator, validationErrors)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldMultipleErrors.Field1", stringValidator1, validationErrors)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldMultipleErrors.Field2", stringValidator2, validationErrors)
	validationErrors = utility.ValidateField(mockLogger, location+"TestValidateFieldMultipleErrors.Field3", stringValidator3, validationErrors)

	expectedLocation1 := location + "TestValidateFieldMultipleErrors.Field1.ValidateField.IsStringLengthInvalid"
	expectedNotification1 := fmt.Sprintf(constants.StringAllowedLength, stringValidator1.MinLength, stringValidator1.MaxLength)
	expectedError1 := domain.NewValidationError(expectedLocation1, stringValidator1.FieldName, constants.FieldRequired, expectedNotification1)

	expectedLocation2 := location + "TestValidateFieldMultipleErrors.Field2.ValidateField.AreStringCharactersInvalid"
	expectedError2 := domain.NewValidationError(expectedLocation2, stringValidator2.FieldName, constants.FieldRequired, constants.StringAllowedCharacters)

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
