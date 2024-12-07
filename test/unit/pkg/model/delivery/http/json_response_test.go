package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	location               = "test.unit.pkg.model.delivery.http."
	field                  = "name"
	validationNotification = "field cannot be empty"
)

func TestNewJSONResponseOnSuccess(t *testing.T) {
	t.Parallel()
	data := "test data"
	response := http.NewJSONResponseOnSuccess(data)

	assert.NoError(t, response.Error, test.DataNilMessage)
	assert.NoError(t, response.Errors, test.DataNilMessage)
	assert.Equal(t, data, response.Data, test.EqualMessage)
	assert.Equal(t, constants.Success, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithHTTPValidationError(t *testing.T) {
	t.Parallel()
	httpBaseError := delivery.NewHTTPValidationError(field, constants.FieldRequired, validationNotification)
	response := http.NewJSONResponseOnFailure(httpBaseError)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.NoError(t, response.Errors, test.ErrorNilMessage)
	assert.IsType(t, delivery.HTTPValidationError{}, response.Error, test.EqualMessage)
	assert.Equal(t, httpBaseError, response.Error, test.EqualMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithMultipleValidationErrors(t *testing.T) {
	t.Parallel()
	errorsList := []error{
		delivery.NewHTTPValidationError(field, constants.FieldRequired, validationNotification),
		delivery.NewHTTPValidationError(field, constants.FieldRequired, validationNotification),
	}
	httpValidationErrors := delivery.NewHTTPValidationErrors(errorsList)
	response := http.NewJSONResponseOnFailure(httpValidationErrors)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.NoError(t, response.Error, test.ErrorNilMessage)
	assert.IsType(t, delivery.HTTPValidationErrors{}, response.Errors, test.EqualMessage)
	assert.Equal(t, httpValidationErrors, response.Errors, test.EqualMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithSingleValidationError(t *testing.T) {
	t.Parallel()
	errorsList := []error{
		delivery.NewHTTPValidationError(field, constants.FieldRequired, validationNotification),
	}
	httpValidationErrors := delivery.NewHTTPValidationErrors(errorsList)
	response := http.NewJSONResponseOnFailure(httpValidationErrors)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.NoError(t, response.Errors, test.ErrorNilMessage)
	assert.IsType(t, delivery.HTTPValidationErrors{}, response.Error, test.EqualMessage)
	assert.Equal(t, httpValidationErrors, response.Error, test.EqualMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithNilError(t *testing.T) {
	t.Parallel()
	response := http.NewJSONResponseOnFailure(nil)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.NoError(t, response.Error, test.DataNilMessage)
	assert.NoError(t, response.Errors, test.ErrorNilMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithEmptyValidationErrorsList(t *testing.T) {
	t.Parallel()
	httpValidationErrors := delivery.NewHTTPValidationErrors([]error{})
	response := http.NewJSONResponseOnFailure(httpValidationErrors)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.NoError(t, response.Error, test.DataNilMessage)
	assert.IsType(t, delivery.HTTPValidationErrors{}, response.Errors, test.EqualMessage)
	assert.Equal(t, httpValidationErrors, response.Errors, test.EqualMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}
