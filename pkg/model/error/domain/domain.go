package domain

import (
	"fmt"
	"strings"
)

// Errors interface represents a collection of errors with a method to get the length of the collection.
type Errors interface {
	error     // Embeds the error interface.
	Len() int // Returns the number of errors in the collection.
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
//
// Parameters:
// - location: The location where the error occurred.
// - notification: The notification message for the error.
//
// Returns:
// - A new instance of BaseError with the provided location and notification.
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
//
// Parameters:
// - errors: A slice of errors to be included in the collection.
//
// Returns:
// - A new instance of BaseErrors with the provided errors.
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
//
// Returns:
// - The number of errors in the collection.
func (baseErrors BaseErrors) Len() int {
	return len(baseErrors.Errors)
}

// InfoMessage represents an informational message, embedding BaseError and implementing the error interface.
type InfoMessage struct {
	BaseError
}

// NewInfoMessage creates a new InfoMessage with the given location and notification message.
//
// Parameters:
// - location: The location where the informational message is relevant.
// - notification: The notification message.
//
// Returns:
// - A new instance of InfoMessage.
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
//
// Parameters:
// - location: The location where the validation error occurred.
// - field: The name of the field that caused the validation error.
// - fieldType: The type of the field that caused the validation error.
// - notification: The notification message for the validation error.
//
// Returns:
// - A new instance of ValidationError.
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
//
// Parameters:
// - errors: A slice of errors to be included in the collection.
//
// Returns:
// - A new instance of ValidationErrors.
func NewValidationErrors(errors []error) ValidationErrors {
	return ValidationErrors{NewBaseErrors(errors)}
}

// AuthorizationError represents an authorization error, embedding BaseError and implementing the error interface.
type AuthorizationError struct {
	BaseError
}

// NewAuthorizationError creates a new AuthorizationError with the given location and notification message.
//
// Parameters:
// - location: The location where the authorization error occurred.
// - notification: The notification message for the authorization error.
//
// Returns:
// - A new instance of AuthorizationError.
func NewAuthorizationError(location, notification string) AuthorizationError {
	return AuthorizationError{NewBaseError(location, notification)}
}

// ItemNotFoundError represents an item not found error, embedding BaseError and adding query-specific details.
type ItemNotFoundError struct {
	BaseError
	Query string // The query that caused the item not found error.
}

// NewItemNotFoundError creates a new ItemNotFoundError with the given details.
//
// Parameters:
// - location: The location where the item was not found.
// - query: The query that caused the item not found error.
// - notification: The notification message for the item not found error.
//
// Returns:
// - A new instance of ItemNotFoundError.
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

// NewInvalidTokenError creates a new InvalidTokenError with the given location and notification message.
//
// Parameters:
// - location: The location where the invalid token error occurred.
// - notification: The notification message for the invalid token error.
//
// Returns:
// - A new instance of InvalidTokenError.
func NewInvalidTokenError(location, notification string) InvalidTokenError {
	return InvalidTokenError{NewBaseError(location, notification)}
}

// TimeExpiredError represents a time expired error, embedding BaseError and implementing the error interface.
type TimeExpiredError struct {
	BaseError
}

// NewTimeExpiredError creates a new TimeExpiredError with the given location and notification message.
//
// Parameters:
// - location: The location where the time expired error occurred.
// - notification: The notification message for the time expired error.
//
// Returns:
// - A new instance of TimeExpiredError.
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
//
// Parameters:
// - location: The location where the pagination error occurred.
// - currentPage: The current page number.
// - totalPages: The total number of pages.
// - notification: The notification message for the pagination error.
//
// Returns:
// - A new instance of PaginationError.
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
//
// Parameters:
// - location: The location where the internal error occurred.
// - notification: The notification message for the internal error.
//
// Returns:
// - A new instance of InternalError.
func NewInternalError(location, notification string) InternalError {
	return InternalError{NewBaseError(location, notification)}
}
