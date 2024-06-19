package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

// ValidationErrorToHttpValidationErrorViewMapper maps a domain ValidationError to an HTTP ValidationError view.
// Parameters:
// - validationError: The domain ValidationError to be mapped.
// Returns:
// - An HttpValidationErrorView populated with the field, fieldType, and notification from the domain error.
func ValidationErrorToHttpValidationErrorViewMapper(validationError domainError.ValidationError) HttpValidationErrorView {
	return HttpValidationErrorView{
		Field:         validationError.Field,
		FieldType:     validationError.FieldType,
		HttpBaseError: NewHttpBaseError(validationError.Notification),
	}
}

// ValidationErrorsToHttpValidationErrorsViewMapper maps a slice of domain ValidationErrors to an HTTP ValidationErrors view.
// Parameters:
// - validationErrors: A slice of domain ValidationErrors to be mapped.
// Returns:
// - An HttpValidationErrorsView populated with the mapped validation errors.
func ValidationErrorsToHttpValidationErrorsViewMapper(validationErrors domainError.ValidationErrors) HttpValidationErrorsView {
	httpValidationErrorsView := make([]error, 0, validationErrors.Len())
	for _, validationError := range validationErrors.Errors {
		validationError, ok := validationError.(domainError.ValidationError)
		if ok {
			// Map the specific validation error to the HTTP view.
			httpValidationErrorView := NewHttpValidationError(validationError.Field, validationError.FieldType, validationError.Notification)
			httpValidationErrorsView = append(httpValidationErrorsView, httpValidationErrorView)
		}
	}

	return NewHttpValidationErrorsView(httpValidationErrorsView)
}

// AuthorizationErrorToHttpAuthorizationErrorViewMapper maps a domain AuthorizationError to an HTTP AuthorizationError view.
// Parameters:
// - authorizationError: The domain AuthorizationError to be mapped.
// Returns:
// - An HttpAuthorizationErrorView populated with the notification from the domain error.
func AuthorizationErrorToHttpAuthorizationErrorViewMapper(authorizationError domainError.AuthorizationError) HttpAuthorizationErrorView {
	return HttpAuthorizationErrorView{
		HttpBaseError: NewHttpBaseError(authorizationError.Notification),
	}
}

// ItemNotFoundErrorToHttpItemNotFoundErrorViewMapper maps a domain ItemNotFoundError to an HTTP ItemNotFoundError view.
// Parameters:
// - ItemNotFoundError: The domain ItemNotFoundError to be mapped.
// Returns:
// - An HttpItemNotFoundErrorView populated with the notification from the domain error.
func ItemNotFoundErrorToHttpItemNotFoundErrorViewMapper(itemNotFoundError domainError.ItemNotFoundError) HttpItemNotFoundErrorView {
	return HttpItemNotFoundErrorView{
		HttpBaseError: NewHttpBaseError(itemNotFoundError.Notification),
	}
}

// PaginationErrorToHttpPaginationErrorViewMapper maps a domain PaginationError to an HTTP PaginationError view.
// Parameters:
// - errorMessage: The domain PaginationError to be mapped.
// Returns:
// - An HttpPaginationErrorView populated with the currentPage, totalPages, and notification from the domain error.
func PaginationErrorToHttpPaginationErrorViewMapper(errorMessage domainError.PaginationError) HttpPaginationErrorView {
	return HttpPaginationErrorView{
		CurrentPage:   errorMessage.CurrentPage,
		TotalPages:    errorMessage.TotalPages,
		HttpBaseError: NewHttpBaseError(errorMessage.Notification),
	}
}

// InternalErrorToHttpInternalErrorViewMapper maps a domain InternalError to an HTTP InternalError view.
// Parameters:
// - internalError: The domain InternalError to be mapped.
// Returns:
// - An HttpInternalErrorView populated with the notification from the domain error.
func InternalErrorToHttpInternalErrorViewMapper(internalError domainError.InternalError) HttpInternalErrorView {
	return HttpInternalErrorView{
		HttpBaseError: NewHttpBaseError(internalError.Notification),
	}
}
