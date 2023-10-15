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
	case ErrorMessage:
		return errorType
	case EntityNotFoundError:
		return NewErrorMessage(constants.EntityNotFoundErrorNotification)
	case PaginationError:
		return errorType
	default:
		return NewErrorMessage(constants.InternalErrorNotification)
	}
}
