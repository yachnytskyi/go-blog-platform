package domain

import (
	"fmt"
	"strings"
)

type Errors interface {
	error
	Len() int
}

type BaseError struct {
	Location     string
	Notification string
}

func (baseError BaseError) Error() string {
	return fmt.Sprintf("location: %s notification: %s", baseError.Location, baseError.Notification)
}

func NewBaseError(location, notification string) BaseError {
	return BaseError{
		Location:     location,
		Notification: notification,
	}
}

type BaseErrors struct {
	Errors []error
}

func NewBaseErrors(errors []error) BaseErrors {
	return BaseErrors{Errors: errors}
}

func (baseErrors BaseErrors) Error() string {
	var result strings.Builder
	result.WriteString("[")
	for i, baseError := range baseErrors.Errors {
		if i > 0 {
			result.WriteString(", ")
		}
		result.WriteString(baseError.Error())
	}
	result.WriteString("]")
	return result.String()
}

func (baseErrors BaseErrors) Len() int {
	return len(baseErrors.Errors)
}

type InfoMessage struct {
	BaseError
}

func NewInfoMessage(location, notification string) InfoMessage {
	return InfoMessage{
		BaseError: NewBaseError(location, notification),
	}
}

type ValidationError struct {
	BaseError
	Field     string
	FieldType string
}

func NewValidationError(location, field, fieldType, notification string) ValidationError {
	return ValidationError{
		BaseError: NewBaseError(location, notification),
		Field:     field,
		FieldType: fieldType,
	}
}

func (validationError ValidationError) Error() string {
	return fmt.Sprintf("location: %s notification: %s field: %s type: %s",
		validationError.Location,
		validationError.Notification,
		validationError.Field,
		validationError.FieldType)
}

type ValidationErrors struct {
	BaseErrors
}

func NewValidationErrors(errors []error) ValidationErrors {
	return ValidationErrors{
		BaseErrors: NewBaseErrors(errors),
	}
}

type AuthorizationError struct {
	BaseError
}

func NewAuthorizationError(location, notification string) AuthorizationError {
	return AuthorizationError{
		BaseError: NewBaseError(location, notification),
	}
}

type ItemNotFoundError struct {
	BaseError
	Query string
}

func NewItemNotFoundError(location, query, notification string) ItemNotFoundError {
	return ItemNotFoundError{
		BaseError: NewBaseError(location, notification),
		Query:     query,
	}
}

func (itemNotFoundError ItemNotFoundError) Error() string {
	return fmt.Sprintf("location: %s notification: %s query: %s",
		itemNotFoundError.Location,
		itemNotFoundError.Notification,
		itemNotFoundError.Query)
}

type InvalidTokenError struct {
	BaseError
}

func NewInvalidTokenError(location, notification string) InvalidTokenError {
	return InvalidTokenError{
		BaseError: NewBaseError(location, notification),
	}
}

type TimeExpiredError struct {
	BaseError
}

func NewTimeExpiredError(location, notification string) TimeExpiredError {
	return TimeExpiredError{
		BaseError: NewBaseError(location, notification),
	}
}

type PaginationError struct {
	BaseError
	CurrentPage string
	TotalPages  string
}

func NewPaginationError(location, currentPage, totalPages, notification string) PaginationError {
	return PaginationError{
		BaseError:   NewBaseError(location, notification),
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}
}

func (paginationError PaginationError) Error() string {
	return fmt.Sprintf("location: %s notification: %s current_page: %s total_pages: %s",
		paginationError.Location,
		paginationError.Notification,
		paginationError.CurrentPage,
		paginationError.TotalPages)
}

type InternalError struct {
	BaseError
}

func NewInternalError(location, notification string) InternalError {
	return InternalError{
		BaseError: NewBaseError(location, notification),
	}
}
