package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type HTTPBaseError struct {
	Location     string `json:"location"`
	Notification string `json:"notification"`
}

func (httpBaseError HTTPBaseError) Error() string {
	return fmt.Sprintf("notification: %s", httpBaseError.Notification)
}

func NewHTTPBaseError(location, notification string) HTTPBaseError {
	return HTTPBaseError{
		Location:     location,
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
	for index, baseError := range httpBaseErrors.Errors {
		if index > 0 {
			result.WriteString(": ")
		}
		result.WriteString(baseError.Error())
	}

	return result.String()
}

func (httpBaseErrors HTTPBaseErrors) Len() int {
	return len(httpBaseErrors.Errors)
}

type HTTPInfoMessage struct {
	HTTPBaseError
}

func NewHTTPInfoMessage(location, notification string) HTTPInfoMessage {
	return HTTPInfoMessage{
		NewHTTPBaseError(location, notification),
	}
}

type HTTPValidationError struct {
	Field     string `json:"field"`
	FieldType string `json:"type"`
	HTTPBaseError
}

func NewHTTPValidationError(location, field, fieldType, notification string) HTTPValidationError {
	return HTTPValidationError{
		Field:         field,
		FieldType:     fieldType,
		HTTPBaseError: NewHTTPBaseError(location, notification),
	}
}

func (httpValidationError HTTPValidationError) Error() string {
	return fmt.Sprintf("%s field: %s type: %s", httpValidationError.HTTPBaseError.Error(), httpValidationError.Field, httpValidationError.FieldType)
}

type HTTPValidationErrors struct {
	HTTPBaseErrors
}

func NewHTTPValidationErrors(errors []error) HTTPValidationErrors {
	return HTTPValidationErrors{NewHTTPBaseErrors(errors)}
}

type HTTPAuthorizationError struct {
	Location string `json:"location"`
	HTTPBaseError
}

func NewHTTPAuthorizationError(location, notification string) HTTPAuthorizationError {
	return HTTPAuthorizationError{
		Location:      location,
		HTTPBaseError: NewHTTPBaseError(location, notification),
	}
}

func (httpAuthorizationError HTTPAuthorizationError) Error() string {
	return fmt.Sprintf("location: %s %s", httpAuthorizationError.Location, httpAuthorizationError.HTTPBaseError.Error())
}

type HTTPItemNotFoundError struct {
	HTTPBaseError
	Query string `json:"query"`
}

func NewHTTPItemNotFoundError(location, notification, query string) HTTPItemNotFoundError {
	return HTTPItemNotFoundError{
		HTTPBaseError: NewHTTPBaseError(location, notification),
		Query:         query,
	}
}

type HTTPInvalidTokenError struct {
	HTTPBaseError
}

func NewHTTPInvalidTokenError(location, notification string) HTTPInvalidTokenError {
	return HTTPInvalidTokenError{NewHTTPBaseError(location, notification)}
}

type HTTPTimeExpiredError struct {
	HTTPBaseError
}

func NewHTTPTimeExpiredError(location, notification string) HTTPTimeExpiredError {
	return HTTPTimeExpiredError{
		HTTPBaseError: NewHTTPBaseError(location, notification),
	}
}

type HTTPPaginationError struct {
	CurrentPage string `json:"current_page"`
	TotalPages  string `json:"total_pages"`
	HTTPBaseError
}

func NewHTTPPaginationError(location, currentPage, totalPages, notification string) HTTPPaginationError {
	return HTTPPaginationError{
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		HTTPBaseError: NewHTTPBaseError(location, notification),
	}
}

func (httpPaginationError HTTPPaginationError) Error() string {
	return fmt.Sprintf("%s current page: %s total pages: %s", httpPaginationError.HTTPBaseError.Error(), httpPaginationError.CurrentPage, httpPaginationError.TotalPages)
}

type HTTPRequestError struct {
	Location    string `json:"location"`
	RequestType string `json:"request_type"`
	HTTPBaseError
}

func NewHTTPRequestError(location, requestType, notification string) HTTPRequestError {
	return HTTPRequestError{
		Location:      location,
		RequestType:   requestType,
		HTTPBaseError: NewHTTPBaseError(location, notification),
	}
}

func (httpRequestError HTTPRequestError) Error() string {
	return fmt.Sprintf("location: %s request type: %s %s", httpRequestError.Location, httpRequestError.RequestType, httpRequestError.HTTPBaseError.Error())
}

type HTTPInternalError struct {
	Location string `json:"location"`
	HTTPBaseError
}

func NewHTTPInternalError(location, notification string) HTTPInternalError {
	return HTTPInternalError{
		Location:      location,
		HTTPBaseError: NewHTTPBaseError(location, notification),
	}
}

func (httpInternalError HTTPInternalError) Error() string {
	return fmt.Sprintf("location: %s %s", httpInternalError.Location, httpInternalError.HTTPBaseError.Error())
}

type HTTPInternalErrors struct {
	HTTPBaseErrors
}

func NewHTTPInternalErrors(errors []error) HTTPInternalErrors {
	return HTTPInternalErrors{NewHTTPBaseErrors(errors)}
}

func (httpValidationErrors HTTPValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpValidationErrors.Errors)
}

func (httpInternalErrors HTTPInternalErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(httpInternalErrors.Errors)
}
