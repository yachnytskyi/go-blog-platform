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
	return fmt.Sprintf("notification: %s", jsonBaseError.Notification)
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
	for index, baseError := range jsonBaseErrors.Errors {
		if index > 0 {
			result.WriteString(": ")
		}
		result.WriteString(baseError.Error())
	}

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
	Field     string `json:"field"`
	FieldType string `json:"type"`
	JSONBaseError
}

func NewJSONValidationError(location, field, fieldType, notification string) JSONValidationError {
	return JSONValidationError{
		Field:         field,
		FieldType:     fieldType,
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

func (jsonValidationError JSONValidationError) Error() string {
	return fmt.Sprintf("%s field: %s type: %s", jsonValidationError.JSONBaseError.Error(), jsonValidationError.Field, jsonValidationError.FieldType)
}

type JSONValidationErrors struct {
	JSONBaseErrors
}

func NewJSONValidationErrors(errors []error) JSONValidationErrors {
	return JSONValidationErrors{NewJSONBaseErrors(errors)}
}

type JSONAuthorizationError struct {
	Location string `json:"location"`
	JSONBaseError
}

func NewJSONAuthorizationError(location, notification string) JSONAuthorizationError {
	return JSONAuthorizationError{
		Location:      location,
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

func (jsonAuthorizationError JSONAuthorizationError) Error() string {
	return fmt.Sprintf("location: %s %s", jsonAuthorizationError.Location, jsonAuthorizationError.JSONBaseError.Error())
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

func (jsonPaginationError JSONPaginationError) Error() string {
	return fmt.Sprintf("%s current page: %s total pages: %s", jsonPaginationError.JSONBaseError.Error(), jsonPaginationError.CurrentPage, jsonPaginationError.TotalPages)
}

type JSONRequestError struct {
	Location    string `json:"location"`
	RequestType string `json:"request_type"`
	JSONBaseError
}

func NewJSONRequestError(location, requestType, notification string) JSONRequestError {
	return JSONRequestError{
		Location:      location,
		RequestType:   requestType,
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

func (jsonRequestError JSONRequestError) Error() string {
	return fmt.Sprintf("location: %s request type: %s %s", jsonRequestError.Location, jsonRequestError.RequestType, jsonRequestError.JSONBaseError.Error())
}

type JSONInternalError struct {
	Location string `json:"location"`
	JSONBaseError
}

func NewJSONInternalError(location, notification string) JSONInternalError {
	return JSONInternalError{
		Location:      location,
		JSONBaseError: NewJSONBaseError(location, notification),
	}
}

func (jsonInternalError JSONInternalError) Error() string {
	return fmt.Sprintf("location: %s %s", jsonInternalError.Location, jsonInternalError.JSONBaseError.Error())
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
