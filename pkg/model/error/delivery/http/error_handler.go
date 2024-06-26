package http

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location = "pkg.model.error.delivery.http." // Constant representing the location of the error handling module.
)

// HandleError takes an error as input and returns an error.
// It performs error type assertions and maps specific domain errors to their corresponding HTTP errors.
//
// Parameters:
// - err: The error to be handled.
//
// Returns:
// - An error mapped to its corresponding HTTP error .
func HandleError(err error) error {
	switch errorType := err.(type) {
	case domainError.ValidationError:
		// Map domain validation error to HTTP validation error .
		return ValidationErrorToHTTPValidationErrorMapper(errorType)
	case domainError.ValidationErrors:
		// Map domain validation errors to HTTP validation errors .
		return ValidationErrorsToHTTPValidationErrorsMapper(errorType)
	case domainError.AuthorizationError:
		// Map domain authorization error to HTTP authorization error .
		return AuthorizationErrorToHTTPAuthorizationErrorMapper(errorType)
	case domainError.ItemNotFoundError:
		// Map domain item not found error to HTTP item not found error .
		return ItemNotFoundErrorToHTTPItemNotFoundErrorMapper(errorType)
	case domainError.PaginationError:
		// Map domain pagination error to HTTP pagination error .
		return PaginationErrorToHTTPPaginationErrorMapper(errorType)
	case HTTPAuthorizationError:
		// Return HTTP authorization error directly.
		return errorType
	case HTTPRequestError:
		// Return HTTP request error directly.
		return errorType
	case HTTPInternalError:
		// Add internal error notification and return HTTP internal error .
		errorType.Notification = constants.InternalErrorNotification
		return errorType
	case HTTPInternalErrors:
		// Return a new HTTP internal error  with location and notification.
		return NewHTTPInternalError(location+"case HTTPInternalErrors", constants.InternalErrorNotification)
	default:
		// Return a new HTTP internal error  for unknown error types with location and notification.
		return NewHTTPInternalError(location+"case default", constants.InternalErrorNotification)
	}
}
