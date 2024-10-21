package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
	mock "github.com/yachnytskyi/golang-mongo-grpc/test/unit/mock/common"
)

const (
	location           = "test.unit.pkg.utility.common."
	originalString     = "test"
	encodedString      = "dGVzdA=="
	emptyString        = ""
	invalidString      = "invalid base64"
	decodeErrorMessage = "illegal base64 data at input byte "
)

func TestDecodeValidBase64String(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	result := utility.Decode(mockLogger, location+"TestDecodeValidBase64String", encodedString)

	assert.False(t, validator.IsError(result.Error), test.NotFailureMessage)
	assert.Equal(t, originalString, result.Data, test.EqualMessage)
}

func TestDecodeEmptyString(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()
	result := utility.Decode(mockLogger, location+"TestDecodeEmptyString", emptyString)

	assert.False(t, validator.IsError(result.Error), test.NotFailureMessage)
	assert.Equal(t, emptyString, result.Data, test.EqualMessage)
}

func TestDecodeInvalidBase64String(t *testing.T) {
	t.Parallel()
	mockLogger := mock.NewMockLogger()

	result := utility.Decode(mockLogger, location+"TestDecodeInvalidBase64String", invalidString)
	expectedLocation := location + "TestDecodeInvalidBase64String.Decode.DecodeString"
	expectedErrorMessage := fmt.Sprintf(test.ExpectedErrorMessageFormat, expectedLocation, decodeErrorMessage+"7")

	assert.True(t, validator.IsError(result.Error), test.FailureMessage)
	assert.IsType(t, domain.InternalError{}, result.Error, test.EqualMessage)
	assert.Equal(t, expectedErrorMessage, result.Error.Error(), test.EqualMessage)
}
