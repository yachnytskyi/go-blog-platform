package http

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/http"
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
		return httpModel.NewJsonResponseOnFailure(NewHttpErrorMessage(config.InternalErrorNotification))
	}
}
