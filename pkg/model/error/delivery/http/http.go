package http

import (
	"encoding/json"
	"fmt"
	"strings"
)

// HttpBaseError represents a base structure for individual errors, containing a notification message.
type HttpBaseError struct {
	Notification string `json:"notification"` // The notification message for the error.
}

// Error implements the error interface for HttpBaseError, returning a formatted error string.
func (httpBaseError HttpBaseError) Error() string {
	return fmt.Sprintf("notification: %s", httpBaseError.Notification)
}

// NewHttpBaseError creates a new HttpBaseError with the given notification message.
func NewHttpBaseError(notification string) HttpBaseError {
	return HttpBaseError{
		Notification: notification,
	}
}

// HttpBaseErrors represents a collection of errors.
type HttpBaseErrors struct {
	Errors []error `json:"errors"` // A slice of errors.
}

// NewHttpBaseErrors creates a new HttpBaseErrors with the given slice of errors.
func NewHttpBaseErrors(errors []error) HttpBaseErrors {
	return HttpBaseErrors{Errors: errors}
}

// Error implements the error interface for HttpBaseErrors, returning a concatenated error string.
func (httpBaseErrors HttpBaseErrors) Error() string {
	var result strings.Builder
	for index, baseError := range httpBaseErrors.Errors {
		if index > 0 {
			result.WriteString(": ")
		}
		result.WriteString(baseError.Error())
	}
	return result.String()
}

// Len returns the number of errors in the HttpBaseErrors collection.
func (httpBaseErrors HttpBaseErrors) Len() int {
	return len(httpBaseErrors.Errors)
}

// HttpValidationErrorView represents a validation error with field-specific details.
type HttpValidationErrorView struct {
	Field     string `json:"field"` // The field that caused the validation error.
	FieldType string `json:"type"`  // The type of the field.
	HttpBaseError
}

// NewHttpValidationError creates a new HttpValidationErrorView with the given details.
func NewHttpValidationError(field, fieldType, notification string) HttpValidationErrorView {
	return HttpValidationErrorView{
		Field:         field,
		FieldType:     fieldType,
		HttpBaseError: NewHttpBaseError(notification),
	}
}

// HttpValidationErrorsView represents a collection of validation errors.
type HttpValidationErrorsView struct {
	HttpBaseErrors
}

// NewHttpValidationErrorsView creates a new HttpValidationErrorsView collection with the given slice of errors.
func NewHttpValidationErrorsView(errors []error) HttpValidationErrorsView {
	return HttpValidationErrorsView{NewHttpBaseErrors(errors)}
}

// HttpAuthorizationErrorView represents an authorization error with additional details.
type HttpAuthorizationErrorView struct {
	Location string `json:"-"` // The location where the error occurred.
	HttpBaseError
}

// NewHttpAuthorizationErrorView creates a new HttpAuthorizationErrorView with the given location and notification message.
func NewHttpAuthorizationErrorView(location, notification string) HttpAuthorizationErrorView {
	return HttpAuthorizationErrorView{
		Location:      location,
		HttpBaseError: NewHttpBaseError(notification),
	}
}

// HttpItemNotFoundErrorView represents an item not found error.
type HttpItemNotFoundErrorView struct {
	HttpBaseError
}

// NewHttpItemNotFoundErrorView creates a new HttpItemNotFoundErrorView with the given notification message.
func NewHttpItemNotFoundErrorView(notification string) HttpItemNotFoundErrorView {
	return HttpItemNotFoundErrorView{
		HttpBaseError: NewHttpBaseError(notification),
	}
}

// HttpPaginationErrorView represents a pagination error with additional details.
type HttpPaginationErrorView struct {
	Location    string `json:"-"`            // The location where the error occurred.
	CurrentPage string `json:"current_page"` // The current page number.
	TotalPages  string `json:"total_pages"`  // The total number of pages.
	HttpBaseError
}

// NewHttpPaginationErrorView creates a new HttpPaginationErrorView with the given details.
func NewHttpPaginationErrorView(location, currentPage, totalPages, notification string) HttpPaginationErrorView {
	return HttpPaginationErrorView{
		Location:      location,
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		HttpBaseError: NewHttpBaseError(notification),
	}
}

// HttpRequestErrorView represents a request error with additional details.
type HttpRequestErrorView struct {
	Location    string `json:"-"`            // The location where the error occurred.
	RequestType string `json:"request_type"` // The type of the request that caused the error.
	HttpBaseError
}

// NewHttpRequestErrorView creates a new HttpRequestErrorView with the given location, request type, and notification message.
func NewHttpRequestErrorView(location, requestType, notification string) HttpRequestErrorView {
	return HttpRequestErrorView{
		Location:      location,
		RequestType:   requestType,
		HttpBaseError: NewHttpBaseError(notification),
	}
}

// HttpInternalErrorView represents an internal error with additional details.
type HttpInternalErrorView struct {
	Location string `json:"-"` // The location where the error occurred.
	HttpBaseError
}

// NewHttpInternalErrorView creates a new HttpInternalErrorView with the given location and notification message.
func NewHttpInternalErrorView(location, notification string) HttpInternalErrorView {
	return HttpInternalErrorView{
		Location:      location,
		HttpBaseError: NewHttpBaseError(notification),
	}
}

// HttpInternalErrorsView represents a collection of internal errors.
type HttpInternalErrorsView struct {
	HttpBaseErrors
}

// NewHttpInternalErrorsView creates a new HttpInternalErrorsView collection with the given slice of errors.
func NewHttpInternalErrorsView(errors []error) HttpInternalErrorsView {
	return HttpInternalErrorsView{NewHttpBaseErrors(errors)}
}

// MarshalJSON customizes the JSON output for HttpValidationErrorsView.
func (httpValidationErrorsView HttpValidationErrorsView) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpValidationErrorsView.Errors)
}

// MarshalJSON customizes the JSON output for HttpInternalErrorsView.
func (httpInternalErrorsView HttpInternalErrorsView) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpInternalErrorsView.Errors)
}
