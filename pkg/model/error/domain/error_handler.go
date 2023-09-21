package domain

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
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
		return NewErrorMessage(config.EntityNotFoundErrorNotification)
	case PaginationError:
		return errorType
	default:
		return NewErrorMessage(config.InternalErrorNotification)
	}
}
