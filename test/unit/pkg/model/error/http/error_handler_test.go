package http_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	location = "test.unit.pkg.model.error.http."

	field        = "Username"
	notification = "Some test notification"
)

func TestHandleErrorValidationErrors(t *testing.T) {
	t.Parallel()

	validationErrors := domain.NewValidationErrors([]error{
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, notification),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, notification),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, notification),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, notification),
	})
	result := http.HandleError(validationErrors)

	httpValidationErrors, ok := result.(http.HTTPValidationErrors)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, httpValidationErrors.Len(), 4, test.EqualMessage)
}

func TestHandleErrorValidationError(t *testing.T) {
	t.Parallel()
	validationError := domain.NewValidationError(location+"TestHandleErrorValidationError", field, constants.FieldRequired, notification)
	result := http.HandleError(validationError)

	httpError, ok := result.(http.HTTPValidationError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, validationError.Field, httpError.Field, test.EqualMessage)
	assert.Equal(t, validationError.FieldType, httpError.FieldType, test.EqualMessage)
	assert.Equal(t, validationError.Notification, httpError.Notification, test.EqualMessage)
}

func TestHandleErrorAuthorizationError(t *testing.T) {
	t.Parallel()
	authorizationError := domain.NewAuthorizationError(location+"TestHandleErrorAuthorizationError", notification)
	result := http.HandleError(authorizationError)

	httpAuthorizationError, ok := result.(http.HTTPAuthorizationError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, authorizationError.Location, httpAuthorizationError.Location, ok, test.EqualMessage)
	assert.Equal(t, authorizationError.Notification, httpAuthorizationError.Notification, ok, test.EqualMessage)
}

func TestHandleErrorItemNotFoundError(t *testing.T) {
	t.Parallel()
	query := "some test query"
	itemNotFoundError := domain.NewItemNotFoundError(location, query, notification)
	result := http.HandleError(itemNotFoundError)

	httpItemNotFoundError, ok := result.(http.HTTPItemNotFoundError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, itemNotFoundError.Notification, httpItemNotFoundError.Notification, test.EqualMessage)
}

func TestHandleErrorInvalidTokenError(t *testing.T) {
	t.Parallel()
	invalidTokenError := domain.NewInvalidTokenError(location+"TestHandleErrorInvalidTokenError", notification)
	result := http.HandleError(invalidTokenError)

	httpError, ok := result.(http.HTTPInvalidTokenError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, invalidTokenError.Notification, httpError.Notification, test.EqualMessage)
}

func TestHandleErrorTimeExpiredError(t *testing.T) {
	t.Parallel()
	timeExpiredError := domain.NewTimeExpiredError(location+"TestHandleErrorTimeExpiredError", notification)
	result := http.HandleError(timeExpiredError)

	httpError, ok := result.(http.HTTPTimeExpiredError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, timeExpiredError.Notification, httpError.Notification, test.EqualMessage)
}

func TestHandleErrorPaginationError(t *testing.T) {
	t.Parallel()
	currentPage := "5"
	totalPages := "10"
	paginationError := domain.NewPaginationError(location+"TestHandleErrorPaginationError", currentPage, totalPages, notification)
	result := http.HandleError(paginationError)

	httpError, ok := result.(http.HTTPPaginationError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, paginationError.CurrentPage, httpError.CurrentPage, test.EqualMessage)
	assert.Equal(t, paginationError.TotalPages, httpError.TotalPages, test.EqualMessage)
	assert.Equal(t, paginationError.Notification, httpError.Notification, test.EqualMessage)
}

func TestHandleErrorInternalError(t *testing.T) {
	t.Parallel()
	internalError := domain.NewInternalError(location+"TestHandleErrorInternalError", notification)
	result := http.HandleError(internalError)

	httpError, ok := result.(http.HTTPInternalError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, internalError.Location, httpError.Location, test.EqualMessage)
	assert.Equal(t, internalError.Notification, httpError.Notification, test.EqualMessage)
}

func TestHandleErrorHTTPInternalError(t *testing.T) {
	t.Parallel()
	httpInternalError := http.NewHTTPInternalError(location+"TestHandleErrorHTTPInternalError", notification)
	result := http.HandleError(httpInternalError)

	httpError, ok := result.(http.HTTPInternalError)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, httpInternalError.Location, httpError.Location, test.EqualMessage)
	assert.Equal(t, constants.InternalErrorNotification, httpError.Notification, test.EqualMessage)
}

func TestHandleErrorHTTPUnknownError(t *testing.T) {
	t.Parallel()
	err := fmt.Errorf("notification: %s", notification)
	result := http.HandleError(err)

	assert.IsType(t, err, result, test.EqualMessage)
	assert.Equal(t, err, result, test.EqualMessage)
}

func TestHandleErrorNil(t *testing.T) {
	t.Parallel()

	result := http.HandleError(nil)
	assert.Nil(t, result, test.ErrorNilMessage, test.ErrorNilMessage)
}
