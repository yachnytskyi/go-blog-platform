package domain

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

// HandleError is a central error handling function that categorizes and provides
// consistent error handling for various error types in the application.
// It returns an error type with a standardized notification message based on
// the specific error type.
func HandleError(err error) error {
	switch errorType := err.(type) {
	case ValidationError:
		return errorType
	case ValidationErrors:
		return errorType
	case AuthorizationError:
		errorType.Notification = constants.AuthorizationErrorNotification
		return errorType
	case EntityNotFoundError:
		errorType.Notification = constants.EntityNotFoundErrorNotification
		return errorType
	case PaginationError:
		errorType.Notification = constants.PaginationErrorNotification
		return errorType
	default:
		internalError := errorType.(InternalError)
		internalError.Notification = constants.InternalErrorNotification
		return internalError
	}
}
