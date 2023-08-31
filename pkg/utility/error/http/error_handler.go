package http

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/model/http"
)

func HandleError(err error) *httpModel.JsonResponse {
	switch errorType := err.(type) {
	case *domainError.ValidationError:
		return httpModel.NewJsonResponseWithError(ValidationErrorToHttpValidationErrorViewMapper(errorType))
	case *domainError.ValidationErrors:
		httpValidationErrors := ValidationErrorsToHttpValidationErrorsViewMapper(errorType)
		return httpModel.NewJsonResponseWithError(httpValidationErrors.HttpValidationErrorsView)
	case *domainError.ErrorMessage:
		return &httpModel.JsonResponse{Error: ErrorMessageToErrorMessageViewMapper(errorType)}
	default:
		var defaultError *domainError.ErrorMessage = new(domainError.ErrorMessage)
		defaultError.Notification = config.InternalErrorNotification
		return &httpModel.JsonResponse{Error: ErrorMessageToErrorMessageViewMapper(defaultError)}
	}
}
