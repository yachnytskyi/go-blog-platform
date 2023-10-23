package http

import (
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func HandleError(err error) httpModel.JsonResponse {
	switch errorType := err.(type) {
	case domainError.ValidationError:
		return httpModel.NewJsonResponseOnFailure(ValidationErrorToHttpValidationErrorViewMapper(errorType))
	case domainError.ValidationErrors:
		httpValidationErrors := ValidationErrorsToHttpValidationErrorsViewMapper(errorType)
		return httpModel.NewJsonResponseOnFailure(httpValidationErrors.HttpValidationErrorsView)
	case domainError.AuthorizationError:
		return httpModel.NewJsonResponseOnFailure(AuthorizationErrorToHttpAuthorizationErrorViewMapper(errorType))
	case domainError.EntityNotFoundError:
		return httpModel.NewJsonResponseOnFailure(EntityNotFoundErrorToHttpEntityNotFoundErrorViewMapper(errorType))
	case domainError.PaginationError:
		return httpModel.NewJsonResponseOnFailure(PaginationErrorToHttpPaginationErrorViewMapper(errorType))
	default:
		internalError := errorType.(domainError.InternalError)
		return httpModel.NewJsonResponseOnFailure(InternalErrorToHttpInternalErrorViewMapper(internalError))
	}
}
