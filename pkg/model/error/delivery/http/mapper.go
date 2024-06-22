package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

// ValidationErrorToHTTPValidationErrorViewMapper maps a domain ValidationError to an HTTP ValidationError view.
// Parameters:
// - validationError: The domain ValidationError to be mapped.
// Returns:
// - An HTTPValidationErrorView populated with the field, fieldType, and notification from the domain error.
func ValidationErrorToHTTPValidationErrorViewMapper(validationError domainError.ValidationError) HTTPValidationErrorView {
	return NewHTTPValidationErrorView(
		validationError.Field,
		validationError.FieldType,
		"",
	)
}

// ValidationErrorsToHTTPValidationErrorsViewMapper maps a slice of domain ValidationErrors to an HTTP ValidationErrors view.
// Parameters:
// - validationErrors: A slice of domain ValidationErrors to be mapped.
// Returns:
// - An HTTPValidationErrorsView populated with the mapped validation errors.
func ValidationErrorsToHTTPValidationErrorsViewMapper(validationErrors domainError.ValidationErrors) HTTPValidationErrorsView {
	httpValidationErrorsView := make([]error, 0, validationErrors.Len())
	for _, validationError := range validationErrors.Errors {
		validationError, ok := validationError.(domainError.ValidationError)
		if ok {
			// Map the specific validation error to the HTTP view.
			httpValidationErrorView := NewHTTPValidationErrorView(
				validationError.Field,
				validationError.FieldType,
				validationError.Notification,
			)
			httpValidationErrorsView = append(httpValidationErrorsView, httpValidationErrorView)
		}
	}

	return NewHTTPValidationErrorsView(httpValidationErrorsView)
}

// AuthorizationErrorToHTTPAuthorizationErrorViewMapper maps a domain AuthorizationError to an HTTP AuthorizationError view.
// Parameters:
// - authorizationError: The domain AuthorizationError to be mapped.
// Returns:
// - An HTTPAuthorizationErrorView populated with the notification from the domain error.
func AuthorizationErrorToHTTPAuthorizationErrorViewMapper(authorizationError domainError.AuthorizationError) HTTPAuthorizationErrorView {
	return NewHTTPAuthorizationErrorView(
		"",
		authorizationError.Notification,
	)
}

// ItemNotFoundErrorToHTTPItemNotFoundErrorViewMapper maps a domain ItemNotFoundError to an HTTP ItemNotFoundError view.
// Parameters:
// - ItemNotFoundError: The domain ItemNotFoundError to be mapped.
// Returns:
// - An HTTPItemNotFoundErrorView populated with the notification from the domain error.
func ItemNotFoundErrorToHTTPItemNotFoundErrorViewMapper(itemNotFoundError domainError.ItemNotFoundError) HTTPItemNotFoundErrorView {
	return NewHTTPItemNotFoundErrorView(
		itemNotFoundError.Notification,
	)
}

// PaginationErrorToHTTPPaginationErrorViewMapper maps a domain PaginationError to an HTTP PaginationError view.
// Parameters:
// - errorMessage: The domain PaginationError to be mapped.
// Returns:
// - An HTTPPaginationErrorView populated with the currentPage, totalPages, and notification from the domain error.
func PaginationErrorToHTTPPaginationErrorViewMapper(errorMessage domainError.PaginationError) HTTPPaginationErrorView {
	return NewHTTPPaginationErrorView(
		errorMessage.CurrentPage,
		errorMessage.TotalPages,
		errorMessage.Notification,
	)
}

// InternalErrorToHTTPInternalErrorViewMapper maps a domain InternalError to an HTTP InternalError view.
// Parameters:
// - internalError: The domain InternalError to be mapped.
// Returns:
// - An HTTPInternalErrorView populated with the notification from the domain error.
func InternalErrorToHTTPInternalErrorViewMapper(internalError domainError.InternalError) HTTPInternalErrorView {
	return NewHTTPInternalErrorView(
		internalError.Location,
		internalError.Notification,
	)
}
