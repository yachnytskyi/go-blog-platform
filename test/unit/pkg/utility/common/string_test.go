package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/mock"
)

const (
	originalString = "test"
	encodedString  = "dGVzdA=="
	emptyString    = ""
	invalidString  = "invalid_base64"

	decodeErrorMessage         = "illegal base64 data at input byte "
	ExpectedErrorMessageFormat = "location: %s notification: %s"
)

// Tests for Encode function

// TestEncodeValidString tests the encoding of a valid string.
func TestEncodeValidString(t *testing.T) {
	t.Parallel()

	encoded := utility.Encode(originalString)
	assert.Equal(t, encodedString, encoded, test.EqualMessage)
}

// TestEncodeEmptyString tests the encoding of an empty string.
func TestEncodeEmptyString(t *testing.T) {
	t.Parallel()

	encoded := utility.Encode(emptyString)
	assert.Equal(t, emptyString, encoded, test.EqualMessage)
}

// Tests for Decode function

// TestDecodeValidBase64String tests the decoding of a valid base64 string.
func TestDecodeValidBase64String(t *testing.T) {
	t.Parallel()

	logger := &mock.MockLogger{}
	result := utility.Decode(logger, location+"TestDecodeValidBase64String", encodedString)

	assert.False(t, validator.IsError(result.Error), test.NotFailureMessage)
	assert.Equal(t, originalString, result.Data, test.EqualMessage)
}

// TestDecodeInvalidBase64String tests the decoding of an invalid base64 string.
func TestDecodeInvalidBase64String(t *testing.T) {
	t.Parallel()

	logger := &mock.MockLogger{}
	expectedLocation := location + "TestDecodeInvalidBase64String.Decode.DecodeString"
	result := utility.Decode(logger, location+"TestDecodeInvalidBase64String", invalidString)

	assert.True(t, validator.IsError(result.Error), test.FailureMessage)
	assert.NotNil(t, result.Error, test.ErrorNilMessage)
	expectedErrorMessage := fmt.Sprintf(test.ExpectedErrorMessageFormat, expectedLocation, decodeErrorMessage+"7")
	assert.Equal(t, expectedErrorMessage, result.Error.Error(), expectedErrorMessage)
}

// TestDecodeInvalidBase64StringDifferentError tests decoding of a differently invalid base64 string.
func TestDecodeInvalidBase64StringDifferentError(t *testing.T) {
	t.Parallel()

	logger := &mock.MockLogger{}
	expectedLocation := location + "TestDecodeInvalidBase64StringDifferentError.Decode.DecodeString"
	invalidString := "!!!" + invalidString + "!!!"
	result := utility.Decode(logger, location+"TestDecodeInvalidBase64StringDifferentError", invalidString)

	assert.True(t, validator.IsError(result.Error), test.FailureMessage)
	assert.NotNil(t, result.Error, test.ErrorNilMessage)
	expectedErrorMessage := fmt.Sprintf(test.ExpectedErrorMessageFormat, expectedLocation, decodeErrorMessage+"0")
	assert.Equal(t, expectedErrorMessage, result.Error.Error(), test.CorrectErrorMessage)
}

// TestDecodeEmptyString tests the decoding of an empty string.
func TestDecodeEmptyString(t *testing.T) {
	t.Parallel()

	logger := &mock.MockLogger{}
	result := utility.Decode(logger, location+"TestDecodeEmptyString", emptyString)

	assert.False(t, validator.IsError(result.Error), test.NotFailureMessage)
	assert.Equal(t, emptyString, result.Data, test.EqualMessage)
}

// Tests for ConvertQueryToString function

// TestConvertQueryToStringStringInput tests converting a string query.
func TestConvertQueryToStringStringInput(t *testing.T) {
	t.Parallel()

	query := originalString
	expected := originalString
	result := utility.ConvertQueryToString(query)
	assert.Equal(t, expected, result, test.EqualMessage)
}

// TestConvertQueryToStringIntegerInput tests converting an integer query.
func TestConvertQueryToStringIntegerInput(t *testing.T) {
	t.Parallel()

	query := 123
	expected := "123"
	result := utility.ConvertQueryToString(query)
	assert.Equal(t, expected, result, test.EqualMessage)
}

// TestConvertQueryToStringStructInput tests converting a struct query.
func TestConvertQueryToStringStructInput(t *testing.T) {
	t.Parallel()

	type testStruct struct {
		Field1 string
		Field2 int
	}

	query := testStruct{"value1", 2}
	expected := "{value1 2}"
	result := utility.ConvertQueryToString(query)
	assert.Equal(t, expected, result, test.EqualMessage)
}

// TestConvertQueryToStringNilInput tests converting a nil query.
func TestConvertQueryToStringNilInput(t *testing.T) {
	t.Parallel()

	var query any = nil
	expected := "<nil>"
	result := utility.ConvertQueryToString(query)
	assert.Equal(t, expected, result, test.EqualMessage)
}

// TestConvertQueryToStringUnexpectedType tests converting an unexpected type query.
func TestConvertQueryToStringUnexpectedType(t *testing.T) {
	t.Parallel()

	query := complex(1, 2)
	expected := "(1+2i)"
	result := utility.ConvertQueryToString(query)
	assert.Equal(t, expected, result, test.EqualMessage)
}
