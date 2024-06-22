package http

import (
	"encoding/json"
	"fmt"
	"strings"
)

// HTTPBaseError represents a base structure for individual errors, containing a notification message.
type HTTPBaseError struct {
	Notification string `json:"notification"` // The notification message for the error.
}

// Error implements the error interface for HTTPBaseError, returning a formatted error string.
func (httpBaseError HTTPBaseError) Error() string {
	return fmt.Sprintf("notification: %s", httpBaseError.Notification)
}

// NewHTTPBaseError creates a new HTTPBaseError with the given notification message.
func NewHTTPBaseError(notification string) HTTPBaseError {
	return HTTPBaseError{
		Notification: notification,
	}
}

// HTTPBaseErrors represents a collection of errors.
type HTTPBaseErrors struct {
	Errors []error `json:"errors"` // A slice of errors.
}

// NewHTTPBaseErrors creates a new HTTPBaseErrors with the given slice of errors.
func NewHTTPBaseErrors(errors []error) HTTPBaseErrors {
	return HTTPBaseErrors{Errors: errors}
}

// Error implements the error interface for HTTPBaseErrors, returning a concatenated error string.
func (httpBaseErrors HTTPBaseErrors) Error() string {
	var result strings.Builder
	for index, baseError := range httpBaseErrors.Errors {
		if index > 0 {
			result.WriteString(": ")
		}
		result.WriteString(baseError.Error())
	}
	return result.String()
}

// Len returns the number of errors in the HTTPBaseErrors collection.
func (httpBaseErrors HTTPBaseErrors) Len() int {
	return len(httpBaseErrors.Errors)
}

// HTTPValidationErrorView represents a validation error with field-specific details.
type HTTPValidationErrorView struct {
	Field     string `json:"field"` // The field that caused the validation error.
	FieldType string `json:"type"`  // The type of the field.
	HTTPBaseError
}

// NewHTTPValidationError creates a new HTTPValidationErrorView with the given details.
func NewHTTPValidationErrorView(field, fieldType, notification string) HTTPValidationErrorView {
	return HTTPValidationErrorView{
		Field:         field,
		FieldType:     fieldType,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPValidationErrorsView represents a collection of validation errors.
type HTTPValidationErrorsView struct {
	HTTPBaseErrors
}

// NewHTTPValidationErrorsView creates a new HTTPValidationErrorsView collection with the given slice of errors.
func NewHTTPValidationErrorsView(errors []error) HTTPValidationErrorsView {
	return HTTPValidationErrorsView{NewHTTPBaseErrors(errors)}
}

// HTTPAuthorizationErrorView represents an authorization error with additional details.
type HTTPAuthorizationErrorView struct {
	Location string `json:"-"` // The location where the error occurred.
	HTTPBaseError
}

// NewHTTPAuthorizationErrorView creates a new HTTPAuthorizationErrorView with the given location and notification message.
func NewHTTPAuthorizationErrorView(location, notification string) HTTPAuthorizationErrorView {
	return HTTPAuthorizationErrorView{
		Location:      location,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPItemNotFoundErrorView represents an item not found error.
type HTTPItemNotFoundErrorView struct {
	HTTPBaseError
}

// NewHTTPItemNotFoundErrorView creates a new HTTPItemNotFoundErrorView with the given notification message.
func NewHTTPItemNotFoundErrorView(notification string) HTTPItemNotFoundErrorView {
	return HTTPItemNotFoundErrorView{
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPPaginationErrorView represents a pagination error with additional details.
type HTTPPaginationErrorView struct {
	CurrentPage string `json:"current_page"` // The current page number.
	TotalPages  string `json:"total_pages"`  // The total number of pages.
	HTTPBaseError
}

// NewHTTPPaginationErrorView creates a new HTTPPaginationErrorView with the given details.
func NewHTTPPaginationErrorView(currentPage, totalPages, notification string) HTTPPaginationErrorView {
	return HTTPPaginationErrorView{
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPRequestErrorView represents a request error with additional details.
type HTTPRequestErrorView struct {
	Location    string `json:"-"`            // The location where the error occurred.
	RequestType string `json:"request_type"` // The type of the request that caused the error.
	HTTPBaseError
}

// NewHTTPRequestErrorView creates a new HTTPRequestErrorView with the given location, request type, and notification message.
func NewHTTPRequestErrorView(location, requestType, notification string) HTTPRequestErrorView {
	return HTTPRequestErrorView{
		Location:      location,
		RequestType:   requestType,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPInternalErrorView represents an internal error with additional details.
type HTTPInternalErrorView struct {
	Location string `json:"-"` // The location where the error occurred.
	HTTPBaseError
}

// NewHTTPInternalErrorView creates a new HTTPInternalErrorView with the given location and notification message.
func NewHTTPInternalErrorView(location, notification string) HTTPInternalErrorView {
	return HTTPInternalErrorView{
		Location:      location,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPInternalErrorsView represents a collection of internal errors.
type HTTPInternalErrorsView struct {
	HTTPBaseErrors
}

// NewHTTPInternalErrorsView creates a new HTTPInternalErrorsView collection with the given slice of errors.
func NewHTTPInternalErrorsView(errors []error) HTTPInternalErrorsView {
	return HTTPInternalErrorsView{NewHTTPBaseErrors(errors)}
}

// MarshalJSON customizes the JSON output for HTTPValidationErrorsView.
func (httpValidationErrorsView HTTPValidationErrorsView) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpValidationErrorsView.Errors)
}

// MarshalJSON customizes the JSON output for HTTPInternalErrorsView.
func (httpInternalErrorsView HTTPInternalErrorsView) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpInternalErrorsView.Errors)
}
