package http

import (
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	location               = "test.unit.pkg.model.delivery.http."
	field                  = "name"
	validationNotification = "field cannot be empty"
)

func TestNewJSONResponseOnSuccess(t *testing.T) {
	data := "test data"
	response := http.NewJSONResponseOnSuccess(data)

	assert.NoError(t, response.Error, test.DataNilMessage)
	assert.NoError(t, response.Errors, test.DataNilMessage)
	assert.Equal(t, data, response.Data, test.EqualMessage)
	assert.Equal(t, constants.Success, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithHTTPValidationError(t *testing.T) {
	httpBaseError := httpError.NewHTTPValidationError(field, constants.FieldRequired, validationNotification)
	response := http.NewJSONResponseOnFailure(httpBaseError)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.IsType(t, httpError.HTTPValidationError{}, response.Error, test.EqualMessage)
	assert.NoError(t, response.Errors, test.ErrorNilMessage)
	assert.Equal(t, httpBaseError, response.Error, test.EqualMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithMultipleValidationErrors(t *testing.T) {
	errorsList := []error{
		httpError.NewHTTPValidationError(field, constants.FieldRequired, validationNotification),
		httpError.NewHTTPValidationError(field, constants.FieldRequired, validationNotification),
	}
	httpValidationErrors := httpError.NewHTTPValidationErrors(errorsList)
	response := http.NewJSONResponseOnFailure(httpValidationErrors)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.IsType(t, httpError.HTTPValidationErrors{}, response.Errors, test.EqualMessage)
	assert.NoError(t, response.Error, test.ErrorNilMessage)
	assert.Equal(t, httpValidationErrors, response.Errors, test.EqualMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithSingleValidationError(t *testing.T) {
	errorsList := []error{
		httpError.NewHTTPValidationError(field, constants.FieldRequired, validationNotification),
	}
	httpValidationErrors := httpError.NewHTTPValidationErrors(errorsList)
	response := http.NewJSONResponseOnFailure(httpValidationErrors)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.IsType(t, httpError.HTTPValidationErrors{}, response.Error, test.EqualMessage)
	assert.NoError(t, response.Errors, test.ErrorNilMessage)
	assert.Equal(t, httpValidationErrors, response.Error, test.EqualMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithNilError(t *testing.T) {
	response := http.NewJSONResponseOnFailure(nil)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.NoError(t, response.Error, test.DataNilMessage)
	assert.NoError(t, response.Errors, test.ErrorNilMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}

func TestNewJSONResponseOnFailureWithEmptyValidationErrorsList(t *testing.T) {
	httpValidationErrors := httpError.NewHTTPValidationErrors([]error{})
	response := http.NewJSONResponseOnFailure(httpValidationErrors)

	assert.Nil(t, response.Data, test.DataNilMessage)
	assert.IsType(t, httpError.HTTPValidationErrors{}, response.Errors, test.EqualMessage)
	assert.NoError(t, response.Error, test.DataNilMessage)
	assert.Equal(t, httpValidationErrors, response.Errors, test.EqualMessage)
	assert.Equal(t, constants.Fail, response.Status, test.EqualMessage)
}
