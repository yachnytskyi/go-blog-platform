package http

import (
	"encoding/json"
	"fmt"
	"strings"
)

type HTTPBaseError struct {
	Notification string `json:"notification"`
}

func (httpBaseError HTTPBaseError) Error() string {
	return fmt.Sprintf("notification: %s", httpBaseError.Notification)
}

func NewHTTPBaseError(notification string) HTTPBaseError {
	return HTTPBaseError{
		Notification: notification,
	}
}

type HTTPBaseErrors struct {
	Errors []error `json:"errors"`
}

func NewHTTPBaseErrors(errors []error) HTTPBaseErrors {
	return HTTPBaseErrors{Errors: errors}
}

func (httpBaseErrors HTTPBaseErrors) Error() string {
	var result strings.Builder
	result.WriteString("[")
	for i, baseError := range httpBaseErrors.Errors {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(baseError.Error())
	}
	result.WriteString("]")
	return result.String()
}

func (httpBaseErrors HTTPBaseErrors) Len() int {
	return len(httpBaseErrors.Errors)
}

type HTTPValidationError struct {
	Field     string `json:"field"`
	FieldType string `json:"type"`
	HTTPBaseError
}

func NewHTTPValidationError(field, fieldType, notification string) HTTPValidationError {
	return HTTPValidationError{
		Field:         field,
		FieldType:     fieldType,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

func (httpValidationError HTTPValidationError) Error() string {
	return fmt.Sprintf("field: %s type: %s notification: %s",
		httpValidationError.Field,
		httpValidationError.FieldType,
		httpValidationError.Notification)
}

type HTTPValidationErrors struct {
	HTTPBaseErrors
}

func NewHTTPValidationErrors(errors []error) HTTPValidationErrors {
	return HTTPValidationErrors{NewHTTPBaseErrors(errors)}
}

type HTTPAuthorizationError struct {
	Location string `json:"-"`
	HTTPBaseError
}

func NewHTTPAuthorizationError(location, notification string) HTTPAuthorizationError {
	return HTTPAuthorizationError{
		Location:      location,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

func (httpAuthorizationError HTTPAuthorizationError) Error() string {
	return fmt.Sprintf("location: %s notification: %s",
		httpAuthorizationError.Location,
		httpAuthorizationError.Notification)
}

type HTTPItemNotFoundError struct {
	HTTPBaseError
}

func NewHTTPItemNotFoundError(notification string) HTTPItemNotFoundError {
	return HTTPItemNotFoundError{
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

func (httpItemNotFoundError HTTPItemNotFoundError) Error() string {
	return fmt.Sprintf("notification: %s", httpItemNotFoundError.Notification)
}

type HTTPInvalidTokenError struct {
	HTTPBaseError
}

func NewHTTPInvalidTokenError(notification string) HTTPInvalidTokenError {
	return HTTPInvalidTokenError{NewHTTPBaseError(notification)}
}

func (httpInvalidTokenError HTTPInvalidTokenError) Error() string {
	return fmt.Sprintf("notification: %s", httpInvalidTokenError.Notification)
}

type HTTPTimeExpiredError struct {
	HTTPBaseError
}

func NewHTTPTimeExpiredError(notification string) HTTPTimeExpiredError {
	return HTTPTimeExpiredError{
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

func (httpTimeExpiredError HTTPTimeExpiredError) Error() string {
	return fmt.Sprintf("notification: %s", httpTimeExpiredError.Notification)
}

type HTTPPaginationError struct {
	CurrentPage string `json:"current_page"`
	TotalPages  string `json:"total_pages"`
	HTTPBaseError
}

func NewHTTPPaginationError(currentPage, totalPages, notification string) HTTPPaginationError {
	return HTTPPaginationError{
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

func (httpPaginationError HTTPPaginationError) Error() string {
	return fmt.Sprintf("current_page: %s total_pages: %s notification: %s",
		httpPaginationError.CurrentPage,
		httpPaginationError.TotalPages,
		httpPaginationError.Notification)
}

type HTTPRequestError struct {
	Location    string `json:"-"`
	RequestType string `json:"request_type"`
	HTTPBaseError
}

func NewHTTPRequestError(location, requestType, notification string) HTTPRequestError {
	return HTTPRequestError{
		Location:      location,
		RequestType:   requestType,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

func (httpRequestError HTTPRequestError) Error() string {
	return fmt.Sprintf("location: %s request_type: %s notification: %s",
		httpRequestError.Location,
		httpRequestError.RequestType,
		httpRequestError.Notification)
}

type HTTPInternalError struct {
	Location string `json:"-"`
	HTTPBaseError
}

func NewHTTPInternalError(location, notification string) HTTPInternalError {
	return HTTPInternalError{
		Location:      location,
		HTTPBaseError: NewHTTPBaseError(notification),
	}
}

func (httpInternalError HTTPInternalError) Error() string {
	return fmt.Sprintf("location: %s notification: %s",
		httpInternalError.Location,
		httpInternalError.Notification)
}

type HTTPInternalErrors struct {
	HTTPBaseErrors
}

func NewHTTPInternalErrors(errors []error) HTTPInternalErrors {
	return HTTPInternalErrors{NewHTTPBaseErrors(errors)}
}

func (HTTPValidationErrors HTTPValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(HTTPValidationErrors.Errors)
}

func (HTTPInternalErrors HTTPInternalErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(HTTPInternalErrors.Errors)
}
