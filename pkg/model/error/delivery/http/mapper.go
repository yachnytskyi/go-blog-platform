package http

import (
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

func ValidationErrorToHTTPValidationErrorMapper(validationError domain.ValidationError) HTTPValidationError {
	return NewHTTPValidationError(
		validationError.Field,
		validationError.FieldType,
		validationError.Notification,
	)
}

func ValidationErrorsToHTTPValidationErrorsMapper(validationErrors domain.ValidationErrors) HTTPValidationErrors {
	httpValidationErrors := make([]error, 0, len(validationErrors.Errors))
	for _, validationError := range validationErrors.Errors {
		validationError, ok := validationError.(domain.ValidationError)
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

func AuthorizationErrorToHTTPAuthorizationErrorMapper(authorizationError domain.AuthorizationError) HTTPAuthorizationError {
	return NewHTTPAuthorizationError(
		authorizationError.Location,
		authorizationError.Notification,
	)
}

func ItemNotFoundErrorToHTTPItemNotFoundErrorMapper(itemNotFoundError domain.ItemNotFoundError) HTTPItemNotFoundError {
	return NewHTTPItemNotFoundError(
		itemNotFoundError.Notification,
	)
}

func InvalidTokenErrorToHTTPIvalidTokenErrorMapper(invalidTokenError domain.InvalidTokenError) HTTPInvalidTokenError {
	return NewHTTPInvalidTokenError(
		invalidTokenError.Notification,
	)
}

func TimeExpiredErrorToHTTPTimeExpiredErrorMapper(timeExpiredError domain.TimeExpiredError) HTTPTimeExpiredError {
	return NewHTTPTimeExpiredError(
		timeExpiredError.Notification,
	)
}

func PaginationErrorToHTTPPaginationErrorMapper(paginationError domain.PaginationError) HTTPPaginationError {
	return NewHTTPPaginationError(
		paginationError.CurrentPage,
		paginationError.TotalPages,
		paginationError.Notification,
	)
}

func InternalErrorToHTTPInternalErrorMapper(internalError domain.InternalError) HTTPInternalError {
	return NewHTTPInternalError(
		internalError.Location,
		internalError.Notification,
	)
}
