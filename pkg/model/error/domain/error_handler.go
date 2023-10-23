package domain

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

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
