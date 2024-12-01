package model

import (
	"encoding/json"
	"fmt"
	"strings"
)

type JSONBaseError struct {
	Location     string `json:"location"`
	Notification string `json:"notification"`
}

func (jsonBaseError JSONBaseError) Error() string {
	return fmt.Sprintf("location: %s notification: %s", jsonBaseError.Location, jsonBaseError.Notification)
}

func NewJSONBaseError(location, notification string) JSONBaseError {
	return JSONBaseError{
		Location:     location,
		Notification: notification,
	}
}

type JSONBaseErrors struct {
	Errors []error `json:"errors"`
}

func NewJSONBaseErrors(errors []error) JSONBaseErrors {
	return JSONBaseErrors{Errors: errors}
}

func (jsonBaseErrors JSONBaseErrors) Error() string {
	var result strings.Builder
	result.WriteString("[")
	for i, baseError := range jsonBaseErrors.Errors {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(baseError.Error())
	}
	
	result.WriteString("]")
	return result.String()
}

func (jsonBaseErrors JSONBaseErrors) Len() int {
	return len(jsonBaseErrors.Errors)
}

type JSONInfoMessage struct {
	JSONBaseError
}

func NewJSONInfoMessage(location, notification string) JSONInfoMessage {
	return JSONInfoMessage{
		NewJSONBaseError(location, notification),
	}
}

type JSONValidationError struct {
	JSONBaseError
	Field     string `json:"field"`
	FieldType string `json:"type"`
}

func NewJSONValidationError(location, field, fieldType, notification string) JSONValidationError {
	return JSONValidationError{
		JSONBaseError: NewJSONBaseError(location, notification),
		Field:         field,
		FieldType:     fieldType,
	}
}

type JSONValidationErrors struct {
	JSONBaseErrors
}

func NewJSONValidationErrors(errors []error) JSONValidationErrors {
	return JSONValidationErrors{NewJSONBaseErrors(errors)}
}

type JSONAuthorizationError struct {
	JSONBaseError
}

func NewJSONAuthorizationError(location, notification string) JSONAuthorizationError {
	return JSONAuthorizationError{
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

type JSONItemNotFoundError struct {
	JSONBaseError
	Query string `json:"query"`
}

func NewJSONItemNotFoundError(location, notification, query string) JSONItemNotFoundError {
	return JSONItemNotFoundError{
		JSONBaseError: NewJSONBaseError(location, notification),
		Query:         query,
	}
}

type JSONInvalidTokenError struct {
	JSONBaseError
}

func NewJSONInvalidTokenError(location, notification string) JSONInvalidTokenError {
	return JSONInvalidTokenError{NewJSONBaseError(location, notification)}
}

type JSONTimeExpiredError struct {
	JSONBaseError
}

func NewJSONTimeExpiredError(location, notification string) JSONTimeExpiredError {
	return JSONTimeExpiredError{
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

type JSONPaginationError struct {
	CurrentPage string `json:"current_page"`
	TotalPages  string `json:"total_pages"`
	JSONBaseError
}

func NewJSONPaginationError(location, currentPage, totalPages, notification string) JSONPaginationError {
	return JSONPaginationError{
		CurrentPage:   currentPage,
		TotalPages:    totalPages,
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

type JSONRequestError struct {
	RequestType string `json:"request_type"`
	JSONBaseError
}

func NewJSONRequestError(location, requestType, notification string) JSONRequestError {
	return JSONRequestError{
		RequestType:   requestType,
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

type JSONInternalError struct {
	JSONBaseError
}

func NewJSONInternalError(location, notification string) JSONInternalError {
	return JSONInternalError{
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

type JSONInternalErrors struct {
	JSONBaseErrors
}

func NewJSONInternalErrors(errors []error) JSONInternalErrors {
	return JSONInternalErrors{NewJSONBaseErrors(errors)}
}

func (JSONValidationErrors JSONValidationErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(JSONValidationErrors.Errors)
}

func (JSONInternalErrors JSONInternalErrors) MarshalJSON() ([]byte, error) {
	return json.Marshal(JSONInternalErrors.Errors)
}
