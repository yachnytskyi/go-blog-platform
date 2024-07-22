package model

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func HandleError(err error) error {
	switch errorType := err.(type) {
	case domainError.InfoMessage:
		return InfoMessageToHTTPInfoMessageMapper(errorType)
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
	case domainError.InternalError:
		return InternalErrorToHTTPInternalErrorMapper(errorType)
	default:
		return errorType
	}
}
