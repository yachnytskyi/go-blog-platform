package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func ValidationErrorToHTTPValidationErrorMapper(validationError domainError.ValidationError) HTTPValidationError {
	return NewHTTPValidationError(
		validationError.Field,
		validationError.FieldType,
		validationError.Notification,
	)
}

func ValidationErrorsToHTTPValidationErrorsMapper(validationErrors domainError.ValidationErrors) HTTPValidationErrors {
	httpValidationErrors := make([]error, 0, validationErrors.Len())
	for _, validationError := range validationErrors.Errors {
		validationError, ok := validationError.(domainError.ValidationError)
		if ok {
			httpValidationError := NewHTTPValidationError(
				validationError.Field,
				validationError.FieldType,
				validationError.Notification,
			)
			httpValidationErrors = append(httpValidationErrors, httpValidationError)
		}
	}

	return NewHTTPValidationErrors(httpValidationErrors)
}

func AuthorizationErrorToHTTPAuthorizationErrorMapper(authorizationError domainError.AuthorizationError) HTTPAuthorizationError {
	return NewHTTPAuthorizationError(
		authorizationError.Location,
		authorizationError.Notification,
	)
}

func ItemNotFoundErrorToHTTPItemNotFoundErrorMapper(itemNotFoundError domainError.ItemNotFoundError) HTTPItemNotFoundError {
	return NewHTTPItemNotFoundError(
		itemNotFoundError.Notification,
	)
}

func InvalidTokenErrorToHTTPIvalidTokenErrorMapper(invalidTokenError domainError.InvalidTokenError) HTTPInvalidTokenError {
	return NewHTTPInvalidTokenError(
		invalidTokenError.Notification,
	)
}

func TimeExpiredErrorToHTTPTimeExpiredErrorMapper(timeExpiredError domainError.TimeExpiredError) HTTPTimeExpiredError {
	return NewHTTPTimeExpiredError(
		timeExpiredError.Notification,
	)
}

func PaginationErrorToHTTPPaginationErrorMapper(paginationError domainError.PaginationError) HTTPPaginationError {
	return NewHTTPPaginationError(
		paginationError.CurrentPage,
		paginationError.TotalPages,
		paginationError.Notification,
	)
}

func InternalErrorToHTTPInternalErrorMapper(internalError domainError.InternalError) HTTPInternalError {
	return NewHTTPInternalError(
		internalError.Location,
		internalError.Notification,
	)
}
