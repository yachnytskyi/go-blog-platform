package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

const (
	originalString = "test"
	encodedString  = "dGVzdA=="
	emptyString    = ""
	invalidString  = "invalid_base64"
	onlyPadding    = "==="

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

// TestEncodeSingleCharacterString tests the encoding of a single character string.
func TestEncodeSingleCharacterString(t *testing.T) {
	t.Parallel()

	encoded := utility.Encode("a")
	assert.Equal(t, "YQ==", encoded, test.EqualMessage)
}

// TestEncodeLongString tests the encoding of a long string.
func TestEncodeLongString(t *testing.T) {
	t.Parallel()

	longString := "a very long string..."
	encoded := utility.Encode(longString)
	decodedResult := utility.Decode(&mock.MockLogger{}, location+"TestEncodeLongString", encoded)
	assert.False(t, validator.IsError(decodedResult.Error), test.NotFailureMessage)
	assert.Equal(t, longString, decodedResult.Data, test.EqualMessage)
}

// TestEncodeBoundaryString tests the encoding of boundary length strings.
func TestEncodeBoundaryString(t *testing.T) {
	t.Parallel()

	boundaryString := "abcd" // Length that results in exact base64 block
	encoded := utility.Encode(boundaryString)
	decodedResult := utility.Decode(&mock.MockLogger{}, location+"TestEncodeBoundaryString", encoded)
	assert.False(t, validator.IsError(decodedResult.Error), test.NotFailureMessage)
	assert.Equal(t, boundaryString, decodedResult.Data, test.EqualMessage)

	boundaryString = "abc" // Length just below base64 padding boundary
	encoded = utility.Encode(boundaryString)
	decodedResult = utility.Decode(&mock.MockLogger{}, location+"TestEncodeBoundaryStringBelow", encoded)
	assert.False(t, validator.IsError(decodedResult.Error), test.NotFailureMessage)
	assert.Equal(t, boundaryString, decodedResult.Data, test.EqualMessage)
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

// TestDecodeOnlyPadding tests the decoding of a string with only padding characters.
func TestDecodeOnlyPadding(t *testing.T) {
	t.Parallel()

	logger := &mock.MockLogger{}
	expectedLocation := location + "TestDecodeOnlyPadding.Decode.DecodeString"
	result := utility.Decode(logger, location+"TestDecodeOnlyPadding", onlyPadding)

	assert.True(t, validator.IsError(result.Error), test.FailureMessage)
	assert.NotNil(t, result.Error, test.ErrorNilMessage)
	expectedErrorMessage := fmt.Sprintf(test.ExpectedErrorMessageFormat, expectedLocation, decodeErrorMessage+"0")
	assert.Equal(t, expectedErrorMessage, result.Error.Error(), test.CorrectErrorMessage)
}

// TestDecodeNilLogger tests
