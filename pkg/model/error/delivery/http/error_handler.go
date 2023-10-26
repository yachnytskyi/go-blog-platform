package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func HandleError(err error) interface{} {
	switch errorType := err.(type) {
	case domainError.ValidationError:
		return ValidationErrorToHttpValidationErrorViewMapper(errorType)
	case domainError.ValidationErrors:
		httpValidationErrors := ValidationErrorsToHttpValidationErrorsViewMapper(errorType)
		return httpValidationErrors.HttpValidationErrorsView
	case domainError.AuthorizationError:
		return AuthorizationErrorToHttpAuthorizationErrorViewMapper(errorType)
	case domainError.EntityNotFoundError:
		return EntityNotFoundErrorToHttpEntityNotFoundErrorViewMapper(errorType)
	case domainError.PaginationError:
		return PaginationErrorToHttpPaginationErrorViewMapper(errorType)
	default:
		internalError := errorType.(domainError.InternalError)
		return InternalErrorToHttpInternalErrorViewMapper(internalError)
	}
}
