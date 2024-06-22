package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

// ValidationErrorToHTTPValidationErrorMapper maps a domain ValidationError to an HTTP ValidationError.
// Parameters:
// - validationError: The domain ValidationError to be mapped.
// Returns:
// - An HTTPValidationError populated with the field, fieldType, and notification from the domain error.
func ValidationErrorToHTTPValidationErrorMapper(validationError domainError.ValidationError) HTTPValidationError {
	return NewHTTPValidationError(
		validationError.Field,
		validationError.FieldType,
		validationError.Location,
	)
}

// ValidationErrorsToHTTPValidationErrorsMapper maps a slice of domain ValidationErrors to an HTTP ValidationErrors.
// Parameters:
// - validationErrors: A slice of domain ValidationErrors to be mapped.
// Returns:
// - An HTTPValidationErrors populated with the mapped validation errors.
func ValidationErrorsToHTTPValidationErrorsMapper(validationErrors domainError.ValidationErrors) HTTPValidationErrors {
	httpValidationErrors := make([]error, 0, validationErrors.Len())
	for _, validationError := range validationErrors.Errors {
		validationError, ok := validationError.(domainError.ValidationError)
		if ok {
			// Map the specific validation error to its corresponding HTTP error.
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

// AuthorizationErrorToHTTPAuthorizationErrorMapper maps a domain AuthorizationError to an HTTP AuthorizationError.
// Parameters:
// - authorizationError: The domain AuthorizationError to be mapped.
// Returns:
// - An HTTPAuthorizationError populated with the notification from the domain error.
func AuthorizationErrorToHTTPAuthorizationErrorMapper(authorizationError domainError.AuthorizationError) HTTPAuthorizationError {
	return NewHTTPAuthorizationError(
		authorizationError.Location,
		authorizationError.Notification,
	)
}

// ItemNotFoundErrorToHTTPItemNotFoundErrorMapper maps a domain ItemNotFoundError to an HTTP ItemNotFoundError.
// Parameters:
// - itemNotFoundError: The domain ItemNotFoundError to be mapped.
// Returns:
// - An HTTPItemNotFoundError populated with the notification from the domain error.
func ItemNotFoundErrorToHTTPItemNotFoundErrorMapper(itemNotFoundError domainError.ItemNotFoundError) HTTPItemNotFoundError {
	return NewHTTPItemNotFoundError(
		itemNotFoundError.Notification,
	)
}

// PaginationErrorToHTTPPaginationErrorMapper maps a domain PaginationError to an HTTP PaginationError.
// Parameters:
// - paginationError: The domain PaginationError to be mapped.
// Returns:
// - An HTTPPaginationError populated with the currentPage, totalPages, and notification from the domain error.
func PaginationErrorToHTTPPaginationErrorMapper(paginationError domainError.PaginationError) HTTPPaginationError {
	return NewHTTPPaginationError(
		paginationError.CurrentPage,
		paginationError.TotalPages,
		paginationError.Notification,
	)
}

// InternalErrorToHTTPInternalErrorMapper maps a domain InternalError to an HTTP InternalError.
// Parameters:
// - internalError: The domain InternalError to be mapped.
// Returns:
// - An HTTPInternalError populated with the notification from the domain error.
func InternalErrorToHTTPInternalErrorMapper(internalError domainError.InternalError) HTTPInternalError {
	return NewHTTPInternalError(
		internalError.Location,
		internalError.Notification,
	)
}
