package utility

import (
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/http_error"
)

func ErrorToErrorViewMapper(err error) error {
	switch errorType := err.(type) {
	case *domainError.ValidationError:
		return httpError.ValidationErrorToHttpValidationErrorViewMapper(errorType)
	case *domainError.ValidationErrors:
		return httpError.ValidationErrorsToHttpValidationErrorsViewMapper(errorType)
	case *domainError.ErrorMessage:
		return httpError.ErrorMessageToErrorMessageViewMapper(errorType)
	default:
		var defaultError *domainError.ErrorMessage = new(domainError.ErrorMessage)
		defaultError.Notification = config.InternalErrorNotification
		return httpError.ErrorMessageToErrorMessageViewMapper(defaultError)
	}
}
