package http

import (
	"fmt"
	"strings"
)

// HttpValidationErrorView represents a validation error in an HTTP request.
type HttpValidationErrorView struct {
	Field        string `json:"field"`        // The field that caused the validation error.
	FieldType    string `json:"type"`         // The type of the field.
	Notification string `json:"notification"` // The notification message for the validation error.
}

// NewHttpValidationError creates a new HttpValidationErrorView with the provided field, fieldType, and notification.
// Parameters:
// - field: The field that caused the validation error.
// - fieldType: The type of the field.
// - notification: The notification message for the validation error.
// Returns:
// - A HttpValidationErrorView struct populated with the given field, fieldType, and notification.
func NewHttpValidationError(field, fieldType, notification string) HttpValidationErrorView {
	return HttpValidationErrorView{
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

// Error implements the error interface for HttpValidationErrorView.
// Returns a formatted string representation of the validation error.
func (httpValidationErrorView HttpValidationErrorView) Error() string {
	return fmt.Sprintf("field: %s type: %s notification: %s", httpValidationErrorView.Field, httpValidationErrorView.FieldType, httpValidationErrorView.Notification)
}

// HttpValidationErrorsView represents a collection of validation errors in an HTTP request.
type HttpValidationErrorsView []HttpValidationErrorView

// NewHttpValidationErrorsView creates a new HttpValidationErrorsView from a slice of HttpValidationErrorView.
// Parameters:
// - httpValidationErrorsView: A slice of HttpValidationErrorView structs.
// Returns:
// - A HttpValidationErrorsView struct populated with the given validation errors.
func NewHttpValidationErrorsView(httpValidationErrorsView []HttpValidationErrorView) HttpValidationErrorsView {
	return HttpValidationErrorsView(httpValidationErrorsView)
}

// Error implements the error interface for HttpValidationErrorsView.
// Returns a concatenated string representation of all validation errors.
func (httpValidationErrorsView HttpValidationErrorsView) Error() string {
	var result strings.Builder
	for _, validationError := range httpValidationErrorsView {
		result.WriteString("field: " + validationError.Field + " type: " + validationError.FieldType + " notification: " + validationError.Notification)
	}
	return result.String()
}

// Len returns the number of validation errors in the HttpValidationErrorsView.
// Returns:
// - An integer representing the number of validation errors.
func (httpValidationErrorsView HttpValidationErrorsView) Len() int {
	return len(httpValidationErrorsView)
}

// HttpAuthorizationErrorView represents an authorization error in an HTTP request.
type HttpAuthorizationErrorView struct {
	Location     string `json:"-"`            // The location where the error occurred.
	Notification string `json:"notification"` // The notification message for the authorization error.
}

// NewHttpAuthorizationErrorView creates a new HttpAuthorizationErrorView with the provided location and notification.
// Parameters:
// - location: The location where the error occurred.
// - notification: The notification message for the authorization error.
// Returns:
// - A HttpAuthorizationErrorView struct populated with the given location and notification.
func NewHttpAuthorizationErrorView(location, notification string) HttpAuthorizationErrorView {
	return HttpAuthorizationErrorView{
		Location:     location,
		Notification: notification,
	}
}

// Error implements the error interface for HttpAuthorizationErrorView.
// Returns a formatted string representation of the authorization error.
func (err HttpAuthorizationErrorView) Error() string {
	return fmt.Sprintf("location: %s notification: %s", err.Location, err.Notification)
}

// HttpEntityNotFoundErrorView represents an entity not found error in an HTTP request.
type HttpEntityNotFoundErrorView struct {
	Notification string `json:"notification"` // The notification message for the entity not found error.
}

// NewHttpEntityNotFoundErrorView creates a new HttpEntityNotFoundErrorView with the provided notification.
// Parameters:
// - notification: The notification message for the entity not found error.
// Returns:
// - A HttpEntityNotFoundErrorView struct populated with the given notification.
func NewHttpEntityNotFoundErrorView(notification string) HttpEntityNotFoundErrorView {
	return HttpEntityNotFoundErrorView{
		Notification: notification,
	}
}

// Error implements the error interface for HttpEntityNotFoundErrorView.
// Returns a formatted string representation of the entity not found error.
func (err HttpEntityNotFoundErrorView) Error() string {
	return fmt.Sprintf("notification: %s", err.Notification)
}

// HttpPaginationErrorView represents a pagination error in an HTTP request.
type HttpPaginationErrorView struct {
	CurrentPage  string `json:"current_page"` // The current page number in the pagination error.
	TotalPages   string `json:"total_pages"`  // The total number of pages in the pagination error.
	Notification string `json:"notification"` // The notification message for the pagination error.
}

// NewHttpPaginationErrorView creates a new HttpPaginationErrorView with the provided currentPage, totalPages, and notification.
// Parameters:
// - currentPage: The current page number in the pagination error.
// - totalPages: The total number of pages in the pagination error.
// - notification: The notification message for the pagination error.
// Returns:
// - A HttpPaginationErrorView struct populated with the given currentPage, totalPages, and notification.
func NewHttpPaginationErrorView(currentPage, totalPages, notification string) HttpPaginationErrorView {
	return HttpPaginationErrorView{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		Notification: notification,
	}
}

// Error implements the error interface for HttpPaginationErrorView.
// Returns a formatted string representation of the pagination error.
func (httpPaginationErrorView HttpPaginationErrorView) Error() string {
	return fmt.Sprintf("current page: %s total pages: %s notification: %s", httpPaginationErrorView.CurrentPage, httpPaginationErrorView.TotalPages, httpPaginationErrorView.Notification)
}

// HttpRequestErrorView represents a request error in an HTTP request.
type HttpRequestErrorView struct {
	RequestType  string `json:"request_type"` // The type of the request that caused the error.
	Notification string `json:"notification"` // The notification message for the request error.
}

// NewHttpRequestErrorView creates a new HttpRequestErrorView with the provided requestType and notification.
// Parameters:
// - requestType: The type of the request that caused the error.
// - notification: The notification message for the request error.
// Returns:
// - A HttpRequestErrorView struct populated with the given requestType and notification.
func NewHttpRequestErrorView(requestType, notification string) HttpRequestErrorView {
	return HttpRequestErrorView{
		RequestType:  requestType,
		Notification: notification,
	}
}

// Error implements the error interface for HttpRequestErrorView.
// Returns a formatted string representation of the request error.
func (httpRequestErrorView HttpRequestErrorView) Error() string {
	return fmt.Sprintf("request type: %s notification: %s", httpRequestErrorView.RequestType, httpRequestErrorView.Notification)
}

// HttpInternalErrorView represents an internal error in an HTTP request.
type HttpInternalErrorView struct {
	Location     string `json:"-"`            // The location where the error occurred.
	Notification string `json:"notification"` // The notification message for the internal error.
}

// NewHttpInternalErrorView creates a new HttpInternalErrorView with the provided location and notification.
// Parameters:
// - location: The location where the error occurred.
// - notification: The notification message for the internal error.
// Returns:
// - A HttpInternalErrorView struct populated with the given location and notification.
func NewHttpInternalErrorView(location, notification string) HttpInternalErrorView {
	return HttpInternalErrorView{
		Location:     location,
		Notification: notification,
	}
}

// Error implements the error interface for HttpInternalErrorView.
// Returns a formatted string representation of the internal error.
func (err HttpInternalErrorView) Error() string {
	return fmt.Sprintf("location: %s notification: %s", err.Location, err.Notification)
}

// HttpInternalErrorsView represents a collection of internal errors in an HTTP request.
type HttpInternalErrorsView []error

// NewHttpInternalErrorsView creates a new HttpInternalErrorsView from a slice of errors.
// Parameters:
// - internalErrors: A slice of error structs.
// Returns:
// - A HttpInternalErrorsView struct populated with the given internal errors.
func NewHttpInternalErrorsView(internalErrors []error) HttpInternalErrorsView {
	return HttpInternalErrorsView(internalErrors)
}

// Error implements the error interface for HttpInternalErrorsView.
// Returns a concatenated string representation of all internal errors.
func (httpInternalErrorsView HttpInternalErrorsView) Error() string {
	var result strings.Builder
	for i, httpInternalView := range httpInternalErrorsView {
		if i > 0 {
			result.WriteString(": ")
		}
		result.WriteString(httpInternalView.Error())
	}

	return result.String()
}

// Len returns the number of internal errors in the HttpInternalErrorsView.
// Returns:
// - An integer representing the number of internal errors.
func (httpInternalErrorsView HttpInternalErrorsView) Len() int {
	return len(httpInternalErrorsView)
}
