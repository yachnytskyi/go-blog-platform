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
//
// Parameters:
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPBaseError with the provided notification message.
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
//
// Parameters:
// - errors: A slice of errors to be included in the collection.
//
// Returns:
// - A new instance of HTTPBaseErrors with the provided errors.
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
//
// Returns:
// - The number of errors in the collection.
func (httpBaseErrors HTTPBaseErrors) Len() int {
	return len(httpBaseErrors.Errors)
}

// HTTPValidationError represents a validation error, embedding HTTPBaseError and adding field-specific details.
type HTTPValidationError struct {
	Field     string `json:"field"` // The field that caused the validation error.
	FieldType string `json:"type"`  // The type of the field.
	HTTPBaseError
}

// NewHTTPValidationError creates a new HTTPValidationError with the given details.
//
// Parameters:
// - field: The field that caused the validation error.
// - fieldType: The type of the field.
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPValidationError with the provided details.
func NewHTTPValidationError(field, fieldType, notification string) HTTPValidationError {
	return HTTPValidationError{
		Field:         field,
		FieldType:     fieldType,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPValidationErrors represents a collection of validation errors.
type HTTPValidationErrors struct {
	HTTPBaseErrors
}

// NewHTTPValidationErrors creates a new HTTPValidationErrors collection with the given slice of errors.
//
// Parameters:
// - errors: A slice of errors to be included in the collection.
//
// Returns:
// - A new instance of HTTPValidationErrors with the provided errors.
func NewHTTPValidationErrors(errors []error) HTTPValidationErrors {
	return HTTPValidationErrors{NewHTTPBaseErrors(errors)}
}

// HTTPAuthorizationError represents an authorization error, embedding HTTPBaseError and adding field-specific details.
type HTTPAuthorizationError struct {
	Location string `json:"-"` // The location where the error occurred.
	HTTPBaseError
}

// NewHTTPAuthorizationError creates a new HTTPAuthorizationError with the given location and notification message.
//
// Parameters:
// - location: The location where the error occurred.
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPAuthorizationError with the provided details.
func NewHTTPAuthorizationError(location, notification string) HTTPAuthorizationError {
	return HTTPAuthorizationError{
		Location:      location,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPItemNotFoundError represents an item not found error, embedding HTTPBaseError and implementing the error interface.
type HTTPItemNotFoundError struct {
	HTTPBaseError
}

// NewHTTPItemNotFoundError creates a new HTTPItemNotFoundError with the given notification message.
//
// Parameters:
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPItemNotFoundError with the provided notification message.
func NewHTTPItemNotFoundError(notification string) HTTPItemNotFoundError {
	return HTTPItemNotFoundError{
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPInvalidTokenError represents an invalid token error, embedding HTTPBaseError and implementing the error interface.
type HTTPInvalidTokenError struct {
	HTTPBaseError
}

// NewHTTPInvalidTokenError creates a new HTTPInvalidTokenError with the given notification message.
//
// Parameters:
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPInvalidTokenError with the provided notification message.
func NewHTTPInvalidTokenError(notification string) HTTPInvalidTokenError {
	return HTTPInvalidTokenError{NewHTTPBaseError(notification)}
}

// HTTPTimeExpiredError represents a time expired error, embedding HTTPBaseError and implementing the error interface.
type HTTPTimeExpiredError struct {
	HTTPBaseError
}

// NewHTTPTimeExpiredError creates a new HTTPTimeExpiredError with the given notification message.
//
// Parameters:
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPTimeExpiredError with the provided notification message.
func NewHTTPTimeExpiredError(notification string) HTTPTimeExpiredError {
	return HTTPTimeExpiredError{
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPPaginationError represents a pagination error, embedding HTTPBaseError and adding field-specific details.
type HTTPPaginationError struct {
	CurrentPage string `json:"current_page"` // The current page number.
	TotalPages  string `json:"total_pages"`  // The total number of pages.
	HTTPBaseError
}

// NewHTTPPaginationError creates a new HTTPPaginationError with the given details.
//
// Parameters:
// - currentPage: The current page number.
// - totalPages: The total number of pages.
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPPaginationError with the provided details.
func NewHTTPPaginationError(currentPage, totalPages, notification string) HTTPPaginationError {
	return HTTPPaginationError{
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPRequestError represents a request error, embedding HTTPBaseError and adding field-specific details.
type HTTPRequestError struct {
	Location    string `json:"-"`            // The location where the error occurred.
	RequestType string `json:"request_type"` // The type of the request that caused the error.
	HTTPBaseError
}

// NewHTTPRequestError creates a new HTTPRequestError with the given location, request type, and notification message.
//
// Parameters:
// - location: The location where the error occurred.
// - requestType: The type of the request that caused the error.
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPRequestError with the provided details.
func NewHTTPRequestError(location, requestType, notification string) HTTPRequestError {
	return HTTPRequestError{
		Location:      location,
		RequestType:   requestType,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPInternalError represents an internal error, embedding HTTPBaseError and adding field-specific details.
type HTTPInternalError struct {
	Location string `json:"-"` // The location where the error occurred.
	HTTPBaseError
}

// NewHTTPInternalError creates a new HTTPInternalError with the given location and notification message.
//
// Parameters:
// - location: The location where the error occurred.
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of HTTPInternalError with the provided details.
func NewHTTPInternalError(location, notification string) HTTPInternalError {
	return HTTPInternalError{
		Location:      location,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

// HTTPInternalErrors represents a collection of internal errors.
type HTTPInternalErrors struct {
	HTTPBaseErrors
}

// NewHTTPInternalErrors creates a new HTTPInternalErrors collection with the given slice of errors.
//
// Parameters:
// - errors: A slice of errors to be included in the collection.
//
// Returns:
// - A new instance of HTTPInternalErrors with the provided errors.
func NewHTTPInternalErrors(errors []error) HTTPInternalErrors {
	return HTTPInternalErrors{NewHTTPBaseErrors(errors)}
}

// MarshalJSON customizes the JSON output for HTTPValidationErrors.
//
// Returns:
// - A JSON-encoded byte slice representing the validation errors.
// - An error if the marshaling fails.
func (httpValidationErrors HTTPValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpValidationErrors.Errors)
}

// MarshalJSON customizes the JSON output for HTTPInternalErrors.
//
// Returns:
// - A JSON-encoded byte slice representing the internal errors.
// - An error if the marshaling fails.
func (httpInternalErrors HTTPInternalErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpInternalErrors.Errors)
}
