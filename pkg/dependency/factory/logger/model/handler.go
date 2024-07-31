package model

import (
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func HandleError(err error) error {
	switch errorType := err.(type) {
	case httpError.HTTPAuthorizationError:
		return HTTPAuthorizationErrorToJSONAuthorizationErrorMapper(errorType)
	case httpError.HTTPRequestError:
		return HTTPRequestErrorToJSONRequestErrorMapper(errorType)
	case httpError.HTTPInternalError:
		return HTTPInternalErrorToJSONInternalErrorMapper(errorType)
	case domainError.InfoMessage:
		return InfoMessageToJSONInfoMessageMapper(errorType)
	case domainError.ValidationError:
		return ValidationErrorToJSONValidationErrorMapper(errorType)
	case domainError.ValidationErrors:
		return ValidationErrorsToJSONValidationErrorsMapper(errorType)
	case domainError.AuthorizationError:
		return AuthorizationErrorToJSONAuthorizationErrorMapper(errorType)
	case domainError.ItemNotFoundError:
		return ItemNotFoundErrorToJSONItemNotFoundErrorMapper(errorType)
	case domainError.InvalidTokenError:
		return InvalidTokenErrorToJSONIvalidTokenErrorMapper(errorType)
	case domainError.TimeExpiredError:
		return TimeExpiredErrorToJSONTimeExpiredErrorMapper(errorType)
	case domainError.PaginationError:
		return PaginationErrorToJSONPaginationErrorMapper(errorType)
	case domainError.InternalError:
		return InternalErrorToJSONInternalErrorMapper(errorType)
	default:
		return errorType
	}
}
