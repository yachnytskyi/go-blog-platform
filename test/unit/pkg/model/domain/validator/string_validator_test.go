package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain/validator"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	location = "test.unit.pkg.model.domain.validator."
)

func TestIsStringLengthInvalidValidLength(t *testing.T) {
	t.Parallel()
	checkedString := "valid"
	minLength := 3
	maxLength := 10
	result := validator.IsStringLengthInvalid(checkedString, minLength, maxLength)

	assert.False(t, result, test.NotFailureMessage)
}

func TestIsStringLengthInvalidTooShort(t *testing.T) {
	t.Parallel()
	checkedString := "no"
	minLength := 3
	maxLength := 10
	result := validator.IsStringLengthInvalid(checkedString, minLength, maxLength)

	assert.True(t, result, test.FailureMessage)
}

func TestIsStringLengthInvalidTooLong(t *testing.T) {
	t.Parallel()
	checkedString := "thisisaverylongstring"
	minLength := 3
	maxLength := 10
	result := validator.IsStringLengthInvalid(checkedString, minLength, maxLength)

	assert.True(t, result, test.FailureMessage)
}

func TestIsStringLengthInvalidExactMinLength(t *testing.T) {
	t.Parallel()
	checkedString := "min"
	minLength := 3
	maxLength := 10
	result := validator.IsStringLengthInvalid(checkedString, minLength, maxLength)

	assert.False(t, result, test.NotFailureMessage)
}

func TestIsStringLengthInvalidExactMaxLength(t *testing.T) {
	t.Parallel()
	checkedString := "maximum"
	minLength := 3
	maxLength := 7
	result := validator.IsStringLengthInvalid(checkedString, minLength, maxLength)

	assert.False(t, result, test.NotFailureMessage)
}

func TestAreStringCharactersInvalidValidCharacters(t *testing.T) {
	t.Parallel()
	checkedString := "Valid123"
	regexPattern := "^[a-zA-Z0-9]+$"
	result := validator.AreStringCharactersInvalid(checkedString, regexPattern)

	assert.False(t, result, test.NotFailureMessage)
}

func TestAreStringCharactersInvalidInvalidCharacters(t *testing.T) {
	t.Parallel()
	checkedString := "Invalid!@#"
	regexPattern := "^[a-zA-Z0-9]+$"

	result := validator.AreStringCharactersInvalid(checkedString, regexPattern)

	assert.True(t, result, test.FailureMessage)
}

func TestAreStringCharactersInvalidEmptyString(t *testing.T) {
	t.Parallel()
	checkedString := ""
	regexPattern := "^[a-zA-Z0-9]+$"
	result := validator.AreStringCharactersInvalid(checkedString, regexPattern)

	assert.True(t, result, test.FailureMessage)
}

func TestAreStringCharactersInvalidComplexRegex(t *testing.T) {
	t.Parallel()
	checkedString := "abc-123_def"
	regexPattern := `^[a-z0-9_\-]+$`
	result := validator.AreStringCharactersInvalid(checkedString, regexPattern)

	assert.False(t, result, test.NotFailureMessage)
}
