package http_error

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
)

func ValidationErrorToHttpValidationErrorViewMapper(validationError *domainError.ValidationError) *HttpValidationErrorView {
	return &HttpValidationErrorView{
		HttpValidationErrorView: &HttpValidationBaseErrorView{validationError.Field, validationError.FieldType, validationError.Notification, "fail"},
	}
}

func ValidationErrorsToHttpValidationErrorsViewMapper(validationErrors *domainError.ValidationErrors) *HttpValidationErrorsView {
	httpValidationErrorsView := make([]*HttpValidationBaseErrorView, 0, len(validationErrors.ValidationErrors))
	for _, validationError := range validationErrors.ValidationErrors {
		httpValidationErrorView := &HttpValidationBaseErrorView{}
		httpValidationErrorView.Field = validationError.Field
		httpValidationErrorView.FieldType = validationError.FieldType
		httpValidationErrorView.Notification = validationError.Notification
		httpValidationErrorsView = append(httpValidationErrorsView, httpValidationErrorView)
	}

	return &HttpValidationErrorsView{
		HttpValidationErrorsView: httpValidationErrorsView, Status: "fail",
	}
}

func ErrorMessageToErrorMessageViewMapper(errorMessage *domainError.ErrorMessage) *HttpErrorMessageView {
	return &HttpErrorMessageView{
		HttpErrorMessageView: &HttpBaseErrorMessageView{errorMessage.Notification, "fail"},
	}
}
