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
	result := common.NewResultOnSuccess(data)

	assert.Equal(t, data, result.Data, test.EqualMessage)
	assert.NoError(t, result.Error, test.ErrorNilMessage)
}

func TestNewResultOnFailureWithError(t *testing.T) {
	t.Parallel()
	err := errors.New(testError)
	result := common.NewResultOnFailure[string](err)

	assert.True(t, validator.IsError(result.Error), test.ErrorNotNilMessage)
	assert.Equal(t, err, result.Error, test.EqualMessage)
	assert.Equal(t, "", result.Data, test.EqualMessage)
}

func TestResultOnFailureWithInternalErrorError(t *testing.T) {
	t.Parallel()
	err := domain.InternalError{}
	result := common.NewResultOnFailure[int](err)

	assert.True(t, validator.IsError(result.Error), test.FailureMessage)
	assert.IsType(t, domain.InternalError{}, result.Error, test.EqualMessage)
	assert.Equal(t, 0, result.Data, test.EqualMessage)
}
