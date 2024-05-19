package http

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location = "pkg.model.error.delivery.http." // Constant representing the location of the error handling module.
)

// HandleError takes an error as input and returns an error.
// It performs error type assertions and maps specific domain errors to their corresponding HTTP error views.
// Parameters:
// - err: The error to be handled.
// Returns:
// - An error mapped to its corresponding HTTP error view.
func HandleError(err error) error {
	switch errorType := err.(type) {
	case domainError.ValidationError:
		// Map domain validation error to HTTP validation error view.
		return ValidationErrorToHttpValidationErrorViewMapper(errorType)
	case domainError.ValidationErrors:
		// Map domain validation errors to HTTP validation errors view.
		return ValidationErrorsToHttpValidationErrorsViewMapper(errorType)
	case domainError.AuthorizationError:
		// Map domain authorization error to HTTP authorization error view.
		return AuthorizationErrorToHttpAuthorizationErrorViewMapper(errorType)
	case domainError.EntityNotFoundError:
		// Map domain entity not found error to HTTP entity not found error view.
		return EntityNotFoundErrorToHttpEntityNotFoundErrorViewMapper(errorType)
	case domainError.PaginationError:
		// Map domain pagination error to HTTP pagination error view.
		return PaginationErrorToHttpPaginationErrorViewMapper(errorType)
	case HttpAuthorizationErrorView:
		// Return HTTP authorization error view directly.
		return errorType
	case HttpRequestErrorView:
		// Return HTTP request error view directly.
		return errorType
	case HttpInternalErrorView:
		// Add internal error notification and return HTTP internal error view.
		errorType.Notification = constants.InternalErrorNotification
		return errorType
	case HttpInternalErrorsView:
		// Return a new HTTP internal error view with location and notification.
		return NewHttpInternalErrorView(location+"case HttpInternalErrorsView", constants.InternalErrorNotification)
	default:
		// Return a new HTTP internal error view for unknown error types with location and notification.
		return NewHttpInternalErrorView(location+"case default", constants.InternalErrorNotification)
	}
}
