package http_error

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
)

func ValidationErrorsToHttpValidationErrorsViewMapper(validationErrors []*domainError.ValidationError) []*HttpValidationErrorView {
	httpValidationErrors := make([]*HttpValidationErrorView, 0)

	for _, validationError := range validationErrors {
		httpValidationError := &HttpValidationErrorView{}
		httpValidationError.Field = validationError.Field
		httpValidationError.FieldType = validationError.FieldType
		httpValidationError.Notification = validationError.Notification
		httpValidationErrors = append(httpValidationErrors, httpValidationError)
	}

	return httpValidationErrors
}

func ValidationErrorToHttpValidationErrorViewMapper(validationError *domainError.ValidationError) *HttpValidationErrorView {
	return &HttpValidationErrorView{
		Field:        validationError.Field,
		FieldType:    validationError.FieldType,
		Notification: validationError.Notification,
	}
}

func ErrorMessageToErrorMessageViewMapper(errorMessage *domainError.ErrorMessage) *HttpErrorMessageView {
	return &HttpErrorMessageView{
		Notification: errorMessage.Notification,
	}
}
