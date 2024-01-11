package http

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

// HandleError takes an error as input and returns error.
// It performs error type assertions and maps specific domain errors to their corresponding HTTP error views.
func HandleError(err error) error {
	switch errorType := err.(type) {
	case domainError.ValidationError:
		return ValidationErrorToHttpValidationErrorViewMapper(errorType)
	case domainError.ValidationErrors:
		return ValidationErrorsToHttpValidationErrorsViewMapper(errorType)
	case domainError.AuthorizationError:
		return AuthorizationErrorToHttpAuthorizationErrorViewMapper(errorType)
	case domainError.EntityNotFoundError:
		return EntityNotFoundErrorToHttpEntityNotFoundErrorViewMapper(errorType)
	case domainError.PaginationError:
		return PaginationErrorToHttpPaginationErrorViewMapper(errorType)
	case HttpAuthorizationErrorView:
		return errorType
	case HttpRequestErrorView:
		return errorType
	case HttpInternalErrorView:
		errorType.Notification = constants.InternalErrorNotification
		return errorType
	default:
		internalError := errorType.(domainError.InternalError)
		internalError.Notification = constants.InternalErrorNotification
		return InternalErrorToHttpInternalErrorViewMapper(internalError)
	}
}
