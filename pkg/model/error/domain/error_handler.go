package domain

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

// HandleError is a central error handling function that categorizes and provides
// consistent error handling for various error types in the application.
// It returns an error type with a standardized notification message based on
// the specific error type.
// Parameters:
// - err: The error to be handled.
// Returns:
// - An error with a standardized notification message.
func HandleError(err error) error {
	switch errorType := err.(type) {
	case ValidationError:
		// Return the validation error as is.
		return errorType
	case ValidationErrors:
		// Return the validation errors as is.
		return errorType
	case AuthorizationError:
		// Set a standardized notification message for authorization errors.
		errorType.Notification = constants.AuthorizationErrorNotification
		return errorType
	case ItemNotFoundError:
		// Set a standardized notification message for item not found errors.
		errorType.Notification = constants.ItemNotFoundErrorNotification
		return errorType
	case PaginationError:
		// Set a standardized notification message for pagination errors.
		errorType.Notification = constants.PaginationErrorNotification
		return errorType
	default:
		// For unhandled error types, assume an InternalError and provide a standard
		// notification message.
		internalError := errorType.(InternalError)
		internalError.Notification = constants.InternalErrorNotification
		return internalError
	}
}
