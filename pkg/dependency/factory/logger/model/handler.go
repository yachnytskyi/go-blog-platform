package model

import (
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func HandleError(err error) error {
	switch errorType := err.(type) {
	case http.HTTPAuthorizationError:
		return HTTPAuthorizationErrorToJSONAuthorizationErrorMapper(errorType)
	case http.HTTPRequestError:
		return HTTPRequestErrorToJSONRequestErrorMapper(errorType)
	case http.HTTPInternalError:
		return HTTPInternalErrorToJSONInternalErrorMapper(errorType)
	case domain.InfoMessage:
		return InfoMessageToJSONInfoMessageMapper(errorType)
	case domain.ValidationError:
		return ValidationErrorToJSONValidationErrorMapper(errorType)
	case domain.ValidationErrors:
		return ValidationErrorsToJSONValidationErrorsMapper(errorType)
	case domain.AuthorizationError:
		return AuthorizationErrorToJSONAuthorizationErrorMapper(errorType)
	case domain.ItemNotFoundError:
		return ItemNotFoundErrorToJSONItemNotFoundErrorMapper(errorType)
	case domain.InvalidTokenError:
		return InvalidTokenErrorToJSONIvalidTokenErrorMapper(errorType)
	case domain.TimeExpiredError:
		return TimeExpiredErrorToJSONTimeExpiredErrorMapper(errorType)
	case domain.PaginationError:
		return PaginationErrorToJSONPaginationErrorMapper(errorType)
	case domain.InternalError:
		return InternalErrorToJSONInternalErrorMapper(errorType)
	default:
		return errorType
	}
}
