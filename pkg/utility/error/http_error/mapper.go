package http_error

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
)

func ValidationErrorToHttpValidationErrorViewMapper(domainValidationErrors []*domainError.ValidationError) []*HttpValidationErrorView {
	httpValidationErrors := make([]*HttpValidationErrorView, 0)

	for _, domaindomainValidationError := range domainValidationErrors {
		httpValidationError := &HttpValidationErrorView{}
		httpValidationError.Field = domaindomainValidationError.Field
		httpValidationError.FieldType = domaindomainValidationError.FieldType
		httpValidationError.Notification = domaindomainValidationError.Notification
		httpValidationErrors = append(httpValidationErrors, httpValidationError)
	}

	return httpValidationErrors
}
