package http

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
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
	case domainError.ErrorMessage:
		return httpModel.NewJsonResponseOnFailure(ErrorMessageToHttpErrorMessageViewMapper(errorType))
	case domainError.PaginationError:
		return httpModel.NewJsonResponseOnFailure(PaginationErrorToHttpPaginationErrorView(errorType))
	default:
		return httpModel.NewJsonResponseOnFailure(NewHttpErrorMessage(constants.InternalErrorNotification))
	}
}
