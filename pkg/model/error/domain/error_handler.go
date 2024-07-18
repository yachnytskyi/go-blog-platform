package domain

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

const (
	location = "pkg.model.error.domain."
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
	case ItemNotFoundError:
		errorType.Notification = constants.ItemNotFoundErrorNotification
		return errorType
	case InvalidTokenError:
		return errorType
	case TimeExpiredError:
		return errorType
	case PaginationError:
		errorType.Notification = constants.PaginationErrorNotification
		return errorType
	case InternalError:
		errorType.Notification = constants.InternalErrorNotification
		return errorType
	default:
		return NewInternalError(location+"HandleError.default", constants.InternalErrorNotification)
	}
}
