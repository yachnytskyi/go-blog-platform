package http_error

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
)

func ValidationErrorToHttpValidationErrorViewMapper(validationError *domainError.ValidationError) *HttpValidationErrorView {
	return &HttpValidationErrorView{
		Field:        validationError.Field,
		FieldType:    validationError.FieldType,
		Notification: validationError.Notification,
	}
}

func ValidationErrorsToHttpValidationErrorsViewMapper(validationErrors *domainError.ValidationErrors) *HttpValidationErrorsView {
	httpValidationErrors := make([]*HttpValidationErrorView, 0, len(validationErrors.ValidationErrors))
	for _, validationError := range validationErrors.ValidationErrors {
		httpValidationErrorView := ValidationErrorToHttpValidationErrorViewMapper(validationError)
		httpValidationErrors = append(httpValidationErrors, httpValidationErrorView)
	}

	return &HttpValidationErrorsView{
		HttpValidationErrorsView: httpValidationErrors,
	}
}

func ErrorMessageToErrorMessageViewMapper(errorMessage *domainError.ErrorMessage) *HttpErrorMessageView {
	return &HttpErrorMessageView{
		Notification: errorMessage.Notification,
	}
}
