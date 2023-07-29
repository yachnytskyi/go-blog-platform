package utility

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/http_error"
)

func ErrorToErrorViewMapper(err error) error {
	switch errorType := err.(type) {
	case *domainError.ValidationError:
		return httpError.ValidationErrorToHttpValidationErrorViewMapper(errorType)
	case *domainError.ErrorMessage:
		return httpError.ErrorMessageToErrorMessageViewMapper(errorType)
	default:
		var defaultError *domainError.ErrorMessage = new(domainError.ErrorMessage)
		defaultError.Notification = InternalErrorNotification
		return httpError.ErrorMessageToErrorMessageViewMapper(defaultError)
	}
}

func ErrorSliceToErrorSliceViewMapper(errors []error) []error {
	for index, errorType := range errors {
		if validationError, ok := errorType.(*domainError.ValidationError); ok {
			validationErrorView := httpError.ValidationErrorToHttpValidationErrorViewMapper(validationError)
			errors[index] = validationErrorView

		} else if errorMessage, ok := errorType.(*domainError.ErrorMessage); ok {
			errorMessageView := httpError.ErrorMessageToErrorMessageViewMapper(errorMessage)
			errors[index] = errorMessageView
			return errors

		} else {
			errorMessageView := httpError.ErrorMessageToErrorMessageViewMapper(errorMessage)
			errorMessageView.Notification = InternalErrorNotification
			errors[index] = errorMessageView
			return errors
		}
	}

	return errors
}
