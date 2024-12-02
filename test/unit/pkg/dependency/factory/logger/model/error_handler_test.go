package model

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger/model"
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	location = "test.unit.pkg.dependency.factory.logger.model."

	field        = "Username"
	notification = "Some test notification"
	requestType  = "GET"
)

func TestHandleErrorHTTPAuthorizationError(t *testing.T) {
	t.Parallel()
	authorizationError := http.NewHTTPAuthorizationError(location+"TestHandleErrorHTTPAuthorizationError", notification)
	result := logger.HandleError(authorizationError)

	jsonAuthorizationError, ok := result.(logger.JSONAuthorizationError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, authorizationError.Location, jsonAuthorizationError.Location, test.EqualMessage)
	assert.Equal(t, authorizationError.Notification, jsonAuthorizationError.Notification, test.EqualMessage)
}

func TestHandleErrorHTTPRequestError(t *testing.T) {
	t.Parallel()
	requestError := http.NewHTTPRequestError(location+"TestHandleErrorHTTPRequestError", requestType, notification)
	result := logger.HandleError(requestError)

	jsonRequestError, ok := result.(logger.JSONRequestError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, requestError.Location, jsonRequestError.Location, test.EqualMessage)
	assert.Equal(t, requestError.RequestType, jsonRequestError.RequestType, test.EqualMessage)
	assert.Equal(t, requestError.Notification, jsonRequestError.Notification, test.EqualMessage)
}

func TestHandleErrorHTTPInternalError(t *testing.T) {
	t.Parallel()
	httpInternalError := http.NewHTTPInternalError(location+"TestHandleErrorHTTPInternalError", notification)
	result := logger.HandleError(httpInternalError)

	jsonInternalError, ok := result.(logger.JSONInternalError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, httpInternalError.Location, jsonInternalError.Location, test.EqualMessage)
	assert.Equal(t, httpInternalError.Notification, jsonInternalError.Notification, test.EqualMessage)
}

func TestHandleErrorInfoMessage(t *testing.T) {
	t.Parallel()
	infoMessage := domain.NewInfoMessage(location+"TestHandleErrorInfoMessage", notification)
	result := logger.HandleError(infoMessage)

	jsonInfoMessage, ok := result.(logger.JSONInfoMessage)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, infoMessage.Location, jsonInfoMessage.Location, test.EqualMessage)
	assert.Equal(t, infoMessage.Notification, jsonInfoMessage.Notification, test.EqualMessage)
}

func TestHandleErrorValidationErrors(t *testing.T) {
	t.Parallel()

	validationErrors := domain.NewValidationErrors([]error{
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, notification),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, notification),
	})
	result := logger.HandleError(validationErrors)

	jsonValidationErrors, ok := result.(logger.JSONValidationErrors)
	assert.True(t, ok, test.EqualMessage)
	assert.Len(t, jsonValidationErrors.Errors, len(validationErrors.Errors), test.EqualMessage)
}

func TestHandleErrorValidationError(t *testing.T) {
	t.Parallel()
	validationError := domain.NewValidationError(location+"TestHandleErrorValidationError", field, constants.FieldRequired, notification)
	result := logger.HandleError(validationError)

	jsonValidationError, ok := result.(logger.JSONValidationError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, validationError.Location, jsonValidationError.Location, test.EqualMessage)
	assert.Equal(t, validationError.Field, jsonValidationError.Field, test.EqualMessage)
	assert.Equal(t, validationError.FieldType, jsonValidationError.FieldType, test.EqualMessage)
	assert.Equal(t, validationError.Notification, jsonValidationError.Notification, test.EqualMessage)
}

func TestHandleErrorAuthorizationError(t *testing.T) {
	t.Parallel()
	authorizationError := domain.NewAuthorizationError(location+"TestHandleErrorAuthorizationError", notification)
	result := logger.HandleError(authorizationError)

	jsonAuthorizationError, ok := result.(logger.JSONAuthorizationError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, authorizationError.Location, jsonAuthorizationError.Location, test.EqualMessage)
	assert.Equal(t, authorizationError.Notification, jsonAuthorizationError.Notification, test.EqualMessage)
}

func TestHandleErrorItemNotFoundError(t *testing.T) {
	t.Parallel()
	query := "some test query"
	itemNotFoundError := domain.NewItemNotFoundError(location, query, notification)
	result := logger.HandleError(itemNotFoundError)

	jsonItemNotFoundError, ok := result.(logger.JSONItemNotFoundError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, itemNotFoundError.Location, jsonItemNotFoundError.Location, test.EqualMessage)
	assert.Equal(t, itemNotFoundError.Query, jsonItemNotFoundError.Query, test.EqualMessage)
	assert.Equal(t, itemNotFoundError.Notification, jsonItemNotFoundError.Notification, test.EqualMessage)
}

func TestHandleErrorInvalidTokenError(t *testing.T) {
	t.Parallel()
	invalidTokenError := domain.NewInvalidTokenError(location+"TestHandleErrorInvalidTokenError", notification)
	result := logger.HandleError(invalidTokenError)

	jsonInvalidTokenError, ok := result.(logger.JSONInvalidTokenError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, invalidTokenError.Location, jsonInvalidTokenError.Location, test.EqualMessage)
	assert.Equal(t, invalidTokenError.Notification, jsonInvalidTokenError.Notification, test.EqualMessage)
}

func TestHandleErrorTimeExpiredError(t *testing.T) {
	t.Parallel()
	timeExpiredError := domain.NewTimeExpiredError(location+"TestHandleErrorTimeExpiredError", notification)
	result := logger.HandleError(timeExpiredError)

	jsonTimeExpiredError, ok := result.(logger.JSONTimeExpiredError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, timeExpiredError.Location, jsonTimeExpiredError.Location, test.EqualMessage)
	assert.Equal(t, timeExpiredError.Notification, jsonTimeExpiredError.Notification, test.EqualMessage)
}

func TestHandleErrorPaginationError(t *testing.T) {
	t.Parallel()
	paginationError := domain.NewPaginationError(location+"TestHandleErrorPaginationError", "5", "10", notification)
	result := logger.HandleError(paginationError)

	jsonPaginationError, ok := result.(logger.JSONPaginationError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, paginationError.Location, jsonPaginationError.Location, test.EqualMessage)
	assert.Equal(t, paginationError.CurrentPage, jsonPaginationError.CurrentPage, test.EqualMessage)
	assert.Equal(t, paginationError.TotalPages, jsonPaginationError.TotalPages, test.EqualMessage)
	assert.Equal(t, paginationError.Notification, jsonPaginationError.Notification, test.EqualMessage)
}

func TestHandleErrorInternalError(t *testing.T) {
	t.Parallel()
	internalError := domain.NewInternalError(location+"TestHandleErrorInternalError", notification)
	result := logger.HandleError(internalError)

	jsonInternalError, ok := result.(logger.JSONInternalError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, internalError.Location, jsonInternalError.Location, test.EqualMessage)
	assert.Equal(t, internalError.Notification, jsonInternalError.Notification, test.EqualMessage)
}

func TestHandleErrorUnknownError(t *testing.T) {
	t.Parallel()
	err := fmt.Errorf("notification: %s", notification)
	result := logger.HandleError(err)

	assert.Equal(t, err, result, test.EqualMessage)
}

func TestHandleErrorNil(t *testing.T) {
	t.Parallel()

	result := logger.HandleError(nil)
	assert.Nil(t, result, test.ErrorNilMessage, test.ErrorNilMessage)
}
