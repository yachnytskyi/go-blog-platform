package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	testError = "test error"
)

func TestNewResultOnSuccess(t *testing.T) {
	t.Parallel()
	data := "test data"
	result := common.NewResultOnSuccess("test data")

	assert.False(t, validator.IsError(result.Error), test.FailureMessage)
	assert.Equal(t, data, result.Data, test.EqualMessage)
	assert.NoError(t, result.Error, test.ErrorNilMessage)
}

func TestNewResultOnFailure(t *testing.T) {
	t.Parallel()
	err := errors.New(testError)
	result := common.NewResultOnFailure[string](err)

	assert.True(t, validator.IsError(result.Error), test.ErrorNotNilMessage)
	assert.Equal(t, err, result.Error, test.EqualMessage)
	assert.Equal(t, "", result.Data, test.EqualMessage)
}

func TestResultIsErrorNoError(t *testing.T) {
	t.Parallel()
	data := 123
	result := common.NewResultOnSuccess(data)

	assert.False(t, validator.IsError(result.Error), test.NotFailureMessage)
}

func TestResultIsErrorWithError(t *testing.T) {
	t.Parallel()
	err := domain.InternalError{}
	result := common.NewResultOnFailure[int](err)

	assert.True(t, validator.IsError(result.Error), test.FailureMessage)
	assert.IsType(t, domain.InternalError{}, result.Error, test.EqualMessage)
}
