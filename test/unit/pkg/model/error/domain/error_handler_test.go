package domain

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	test "github.com/yachnytskyi/golang-mongo-grpc/test"
)

const (
	location     = "test.unit.pkg.model.error.domain."
	field        = "Username"
	notification = "Test error notification"
)

func TestHandleErrorValidationErrors(t *testing.T) {
	t.Parallel()

	validationErrors := domain.NewValidationErrors([]error{
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, constants.StringAllowedCharacters),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, constants.StringAllowedCharacters),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, constants.StringAllowedCharacters),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, constants.StringAllowedCharacters),
	})
	result := domain.HandleError(validationErrors)

	resultValidationErrors, ok := result.(domain.ValidationErrors)
	assert.True(t, ok, test.EqualMessage)
	assert.Equal(t, validationErrors.Len(), 4, test.EqualMessage)
	assert.Equal(t, validationErrors, resultValidationErrors, test.EqualMessage)
}

func TestHandleErrorValidationError(t *testing.T) {
	t.Parallel()
	validationError := domain.NewValidationError(location+"TestHandleErrorValidationError", field, constants.FieldRequired, notification)
	result := domain.HandleError(validationError)

	assert.IsType(t, domain.ValidationError{}, result, test.EqualMessage)
	assert.Equal(t, validationError, result, test.EqualMessage)
}

func TestHandleErrorAuthorizationError(t *testing.T) {
	t.Parallel()
	authorizationError := domain.NewAuthorizationError(location+"TestHandleErrorAuthorizationError", notification)
	result := domain.HandleError(authorizationError)

	assert.IsType(t, domain.AuthorizationError{}, result, test.EqualMessage)
	assert.Equal(t, authorizationError.Location, result.(domain.AuthorizationError).Location, test.EqualMessage)
	assert.Equal(t, constants.AuthorizationErrorNotification, result.(domain.AuthorizationError).Notification, test.EqualMessage)
}

func TestHandleErrorItemNotFoundError(t *testing.T) {
	t.Parallel()
	query := "some test query"
	itemNotFoundError := domain.NewItemNotFoundError(location+"TestHandleErrorItemNotFoundError", query, notification)
	result := domain.HandleError(itemNotFoundError)

	assert.IsType(t, domain.ItemNotFoundError{}, result, test.EqualMessage)
	assert.Equal(t, itemNotFoundError.Location, result.(domain.ItemNotFoundError).Location, test.EqualMessage)
	assert.Equal(t, constants.ItemNotFoundErrorNotification, result.(domain.ItemNotFoundError).Notification, test.EqualMessage)
}

func TestHandleErrorPaginationError(t *testing.T) {
	t.Parallel()
	currentPage := "52"
	totalPages := "100"
	notification := "Pagination error"
	paginationError := domain.NewPaginationError(location+"TestHandleErrorPaginationError", currentPage, totalPages, notification)
	result := domain.HandleError(paginationError)

	assert.IsType(t, domain.PaginationError{}, result, test.EqualMessage)
	assert.Equal(t, paginationError.Location, result.(domain.PaginationError).Location, test.EqualMessage)
	assert.Equal(t, constants.PaginationErrorNotification, result.(domain.PaginationError).Notification, test.EqualMessage)
}

func TestHandleErrorInternalError(t *testing.T) {
	t.Parallel()
	internalError := domain.NewInternalError(location+"TestHandleErrorInternalError", notification)
	result := domain.HandleError(internalError)

	assert.IsType(t, domain.InternalError{}, result, test.EqualMessage)
	assert.Equal(t, internalError.Location, result.(domain.InternalError).Location, test.EqualMessage)
	assert.Equal(t, constants.InternalErrorNotification, result.(domain.InternalError).Notification, test.EqualMessage)
}

func TestHandleErrorUnknownError(t *testing.T) {
	t.Parallel()
	err := fmt.Errorf("notification: %s", notification)
	result := domain.HandleError(err)

	assert.IsType(t, err, result, test.EqualMessage)
	assert.Equal(t, err, result, test.EqualMessage)
}

func TestHandleErrorNil(t *testing.T) {
	t.Parallel()

	result := domain.HandleError(nil)
	assert.Nil(t, result, test.ErrorNilMessage, test.ErrorNilMessage)
}
