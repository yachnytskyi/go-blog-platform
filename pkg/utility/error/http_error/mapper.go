package http_error

import domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"

func ValidationErrorToHttpValidationErrorMapper(domainValidationErrors []*domainError.ValidationError) []*HttpValidationError {
	httpValidationErrors := make([]*HttpValidationError, 0)

	for _, domaindomainValidationError := range domainValidationErrors {
		httpValidationError := &HttpValidationError{}
		httpValidationError.Field = domaindomainValidationError.Field
		httpValidationError.FieldType = domaindomainValidationError.FieldType
		httpValidationError.Notification = domaindomainValidationError.Notification
		httpValidationErrors = append(httpValidationErrors, httpValidationError)
	}

	return httpValidationErrors
}
