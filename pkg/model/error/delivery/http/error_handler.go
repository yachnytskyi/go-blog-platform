package http

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func HandleError(err error) error {
	switch errorType := err.(type) {
	case domain.ValidationError:
		return ValidationErrorToHTTPValidationErrorMapper(errorType)
	case domain.ValidationErrors:
		return ValidationErrorsToHTTPValidationErrorsMapper(errorType)
	case domain.AuthorizationError:
		return AuthorizationErrorToHTTPAuthorizationErrorMapper(errorType)
	case domain.ItemNotFoundError:
		return ItemNotFoundErrorToHTTPItemNotFoundErrorMapper(errorType)
	case domain.InvalidTokenError:
		return InvalidTokenErrorToHTTPIvalidTokenErrorMapper(errorType)
	case domain.TimeExpiredError:
		return TimeExpiredErrorToHTTPTimeExpiredErrorMapper(errorType)
	case domain.PaginationError:
		return PaginationErrorToHTTPPaginationErrorMapper(errorType)
	case domain.InternalError:
		return InternalErrorToHTTPInternalErrorMapper(errorType)
	case HTTPInternalError:
		errorType.Notification = constants.InternalErrorNotification
		return errorType
	default:
		return errorType
	}
}
