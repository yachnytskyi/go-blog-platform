package http

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location = "pkg.model.error.delivery.http."
)

func HandleError(err error) error {
	switch errorType := err.(type) {
	case domainError.ValidationError:
		return ValidationErrorToHTTPValidationErrorMapper(errorType)
	case domainError.ValidationErrors:
		return ValidationErrorsToHTTPValidationErrorsMapper(errorType)
	case domainError.AuthorizationError:
		return AuthorizationErrorToHTTPAuthorizationErrorMapper(errorType)
	case domainError.ItemNotFoundError:
		return ItemNotFoundErrorToHTTPItemNotFoundErrorMapper(errorType)
	case domainError.InvalidTokenError:
		return InvalidTokenErrorToHTTPIvalidTokenErrorMapper(errorType)
	case domainError.TimeExpiredError:
		return TimeExpiredErrorToHTTPTimeExpiredErrorMapper(errorType)
	case domainError.PaginationError:
		return PaginationErrorToHTTPPaginationErrorMapper(errorType)
	case HTTPAuthorizationError:
		return errorType
	case HTTPRequestError:
		return errorType
	case HTTPInternalError:
		errorType.Notification = constants.InternalErrorNotification
		return errorType
	case HTTPInternalErrors:
		return errorType
	default:
		return NewHTTPInternalError(location+"HandleError.default", constants.InternalErrorNotification)
	}
}
