package domain

import (
	"fmt"
	"strings"
)

// Errors interface represents a collection of errors with a method to get the length of the collection.
type Errors interface {
	error
	Len() int
}

// BaseError represents a base structure for individual errors, containing a location and a notification message.
type BaseError struct {
	Location     string // The location where the error occurred.
	Notification string // The error message or notification.
}

// Error implements the error interface for BaseError, returning a formatted error string.
func (baseError BaseError) Error() string {
	return fmt.Sprintf("location: %s notification: %s", baseError.Location, baseError.Notification)
}

// NewBaseError creates a new BaseError with the given location and notification message.
func NewBaseError(location, notification string) BaseError {
	return BaseError{
		Location:     location,
		Notification: notification,
	}
}

// BaseErrors represents a collection of errors.
type BaseErrors struct {
	Errors []error // A slice of errors.
}

// NewBaseErrors creates a new BaseErrors with the given slice of errors.
func NewBaseErrors(errors []error) BaseErrors {
	return BaseErrors{Errors: errors}
}

// Error implements the error interface for BaseErrors, returning a concatenated error string.
func (baseErrors BaseErrors) Error() string {
	var result strings.Builder
	for index, baseError := range baseErrors.Errors {
		if index > 0 {
			result.WriteString(": ")
		}
		result.WriteString(baseError.Error())
	}

	return result.String()
}

// Len returns the number of errors in the BaseErrors collection.
func (baseErrors BaseErrors) Len() int {
	return len(baseErrors.Errors)
}

// InfoMessage represents an informational message, embedding BaseError and implementing the error interface.
type InfoMessage struct {
	BaseError
}

// NewInfoMessage creates a new InfoMessage with the given location and notification message.
func NewInfoMessage(location, notification string) InfoMessage {
	return InfoMessage{NewBaseError(location, notification)}
}

// ValidationError represents a validation error, embedding BaseError and adding field-specific details.
type ValidationError struct {
	BaseError
	Field     string // The name of the field that caused the validation error.
	FieldType string // The type of the field that caused the validation error.
}

// NewValidationError creates a new ValidationError with the given details.
func NewValidationError(location, field, fieldType, notification string) ValidationError {
	return ValidationError{
		BaseError: NewBaseError(location, notification),
		Field:     field,
		FieldType: fieldType,
	}
}

// ValidationErrors represents a collection of validation errors, embedding BaseErrors and implementing the Errors interface.
type ValidationErrors struct {
	BaseErrors
}

// NewValidationErrors creates a new ValidationErrors collection with the given slice of errors.
func NewValidationErrors(errors []error) ValidationErrors {
	return ValidationErrors{NewBaseErrors(errors)}
}

// AuthorizationError represents an authorization error, embedding BaseError and implementing the error interface.
type AuthorizationError struct {
	BaseError
}

// NewAuthorizationError creates a new AuthorizationError with the given location and notification message.
func NewAuthorizationError(location, notification string) AuthorizationError {
	return AuthorizationError{NewBaseError(location, notification)}
}

// ItemNotFoundError represents an item not found error, embedding BaseError and adding query-specific details.
type ItemNotFoundError struct {
	BaseError
	Query string // The query that caused the item not found error.
}

// NewItemNotFoundError creates a new ItemNotFoundError with the given details.
func NewItemNotFoundError(location, query, notification string) ItemNotFoundError {
	return ItemNotFoundError{
		BaseError: NewBaseError(location, notification),
		Query:     query,
	}
}

// InvalidTokenError represents an invalid token error, embedding BaseError and implementing the error interface.
type InvalidTokenError struct {
	BaseError
}

// InvalidTokenError creates a new InvalidTokenError with the given location and notification message.
func NewInvalidTokenError(location, notification string) InvalidTokenError {
	return InvalidTokenError{NewBaseError(location, notification)}
}

// TimeExpiredError represents a time expired error, embedding BaseError and implementing the error interface.
type TimeExpiredError struct {
	BaseError
}

// TimeExpiredError creates a new TimeExpiredError with the given location and notification message.
func NewTimeExpiredError(location, notification string) TimeExpiredError {
	return TimeExpiredError{NewBaseError(location, notification)}
}

// PaginationError represents a pagination error, embedding BaseError and adding pagination-specific details.
type PaginationError struct {
	BaseError
	CurrentPage string // The current page number.
	TotalPages  string // The total number of pages.
}

// NewPaginationError creates a new PaginationError with the given details.
func NewPaginationError(location, currentPage, totalPages, notification string) PaginationError {
	return PaginationError{
		BaseError:   NewBaseError(location, notification),
		CurrentPage: currentPage,
		TotalPages:  totalPages,
	}
}

// InternalError represents an internal error, embedding BaseError and implementing the error interface.
type InternalError struct {
	BaseError
}

// NewInternalError creates a new InternalError with the given location and notification message.
func NewInternalError(location, notification string) InternalError {
	return InternalError{NewBaseError(location, notification)}
}
