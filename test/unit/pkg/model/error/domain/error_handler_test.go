package domain

import (
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
	validationError := domain.NewValidationErrors([]error{
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, constants.StringAllowedCharacters),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, constants.StringAllowedCharacters),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, constants.StringAllowedCharacters),
		domain.NewValidationError(location+"TestHandleErrorValidationErrors", field, constants.FieldRequired, constants.StringAllowedCharacters),
	})
	result := domain.HandleError(validationError)

	assert.Len(t, validationError.Errors, 4, test.EqualMessage)
	assert.IsType(t, domain.ValidationErrors{}, result, test.EqualMessage)
	assert.Equal(t, validationError, result)

}

func TestHandleErrorAuthorizationError(t *testing.T) {
	t.Parallel()
	authorizationError := domain.NewAuthorizationError(location, notification)
	result := domain.HandleError(authorizationError)

	assert.IsType(t, domain.AuthorizationError{}, result)
	assert.Equal(t, constants.AuthorizationErrorNotification, result.(domain.AuthorizationError).Notification)
}
