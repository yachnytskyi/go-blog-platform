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

// InfoMessage represents an informational message with details about the location and a notification message.
type InfoMessage struct {
	Location     string // The location where the message originated.
	Notification string // The notification message.
}

// NewInfoMessage creates a new InfoMessage with the provided details.
// Parameters:
// - location: The location where the message originated.
// - notification: The notification message.
// Returns:
// - An InfoMessage populated with the given details.
func NewInfoMessage(location, notification string) InfoMessage {
	return InfoMessage{
		Location:     location,
		Notification: notification,
	}
}

// Error implements the error interface for InfoMessage.
// Returns a formatted string representation of the informational message.
func (msg InfoMessage) Error() string {
	return fmt.Sprintf("location: %s notification: %s", msg.Location, msg.Notification)
}

// ValidationError represents a validation error with details about the location, field, field type, and a notification message.
type ValidationError struct {
	Location     string // The location where the error occurred.
	Field        string // The field that caused the validation error.
	FieldType    string // The type of the field.
	Notification string // The notification message for the validation error.
}

// NewValidationError creates a new ValidationError with the provided details.
// Parameters:
// - location: The location where the error occurred.
// - field: The field that caused the validation error.
// - fieldType: The type of the field.
// - notification: The notification message for the validation error.
// Returns:
// - A ValidationError populated with the given details.
func NewValidationError(location, field, fieldType, notification string) ValidationError {
	return ValidationError{
		Location:     location,
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

// Error implements the error interface for ValidationError.
// Returns a formatted string representation of the validation error.
func (err ValidationError) Error() string {
	return fmt.Sprintf("location: %s field: %s type: %s notification: %s", err.Location, err.Field, err.FieldType, err.Notification)
}

// ValidationErrors represents a collection of validation errors.
type ValidationErrors []error

// NewValidationErrors creates a new ValidationErrors collection from a slice of errors.
// Parameters:
// - validationErrors: A slice of errors representing validation errors.
// Returns:
// - A ValidationErrors collection populated with the given errors.
func NewValidationErrors(validationErrors []error) ValidationErrors {
	return ValidationErrors(validationErrors)
}

// Error implements the error interface for ValidationErrors.
// Returns a concatenated string representation of all validation errors.
func (validationErrors ValidationErrors) Error() string {
	var result strings.Builder
	for i, validationError := range validationErrors {
		if i > 0 {
			result.WriteString(": ")
		}
		result.WriteString(validationError.Error())
	}
	return result.String()
}

// Len returns the number of validation errors in the ValidationErrors collection.
// Returns:
// - An integer representing the number of validation errors.
func (validationErrors ValidationErrors) Len() int {
	return len(validationErrors)
}

// AuthorizationError represents an authorization error with details about the location and a notification message.
type AuthorizationError struct {
	Location     string // The location where the error occurred.
	Notification string // The notification message for the authorization error.
}

// NewAuthorizationError creates a new AuthorizationError with the provided details.
// Parameters:
// - location: The location where the error occurred.
// - notification: The notification message for the authorization error.
// Returns:
// - An AuthorizationError populated with the given details.
func NewAuthorizationError(location, notification string) AuthorizationError {
	return AuthorizationError{
		Location:     location,
		Notification: notification,
	}
}

// Error implements the error interface for AuthorizationError.
// Returns a formatted string representation of the authorization error.
func (err AuthorizationError) Error() string {
	return fmt.Sprintf("location: %s notification: %s", err.Location, err.Notification)
}

// ItemNotFoundError represents an item not found error with details about the location, query, and a notification message.
type ItemNotFoundError struct {
	Location     string // The location where the error occurred.
	Query        string // The query that caused the error.
	Notification string // The notification message for the item not found error.
}

// NewItemNotFoundError creates a new ItemNotFoundError with the provided details.
// Parameters:
// - location: The location where the error occurred.
// - query: The query that caused the error.
// - notification: The notification message for the item not found error.
// Returns:
// - An ItemNotFoundError populated with the given details.
func NewItemNotFoundError(location, query, notification string) ItemNotFoundError {
	return ItemNotFoundError{
		Location:     location,
		Query:        query,
		Notification: notification,
	}
}

// Error implements the error interface for ItemNotFoundError.
// Returns a formatted string representation of the item not found error.
func (err ItemNotFoundError) Error() string {
	return fmt.Sprintf("location: %s query: %s notification: %s", err.Location, err.Query, err.Notification)
}

// PaginationError represents a pagination error with details about the current page, total pages, and a notification message.
type PaginationError struct {
	CurrentPage  string // The current page in the pagination.
	TotalPages   string // The total number of pages.
	Notification string // The notification message for the pagination error.
}

// NewPaginationError creates a new PaginationError with the provided details.
// Parameters:
// - currentPage: The current page in the pagination.
// - totalPages: The total number of pages.
// - notification: The notification message for the pagination error.
// Returns:
// - A PaginationError populated with the given details.
func NewPaginationError(currentPage, totalPages, notification string) PaginationError {
	return PaginationError{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		Notification: notification,
	}
}

// Error implements the error interface for PaginationError.
// Returns a formatted string representation of the pagination error.
func (err PaginationError) Error() string {
	return fmt.Sprintf("page: %s total: %s notification: %s", err.CurrentPage, err.TotalPages, err.Notification)
}

// InternalError represents an internal error with details about the location and a notification message.
type InternalError struct {
	Location     string // The location where the error occurred.
	Notification string // The notification message for the internal error.
}

// NewInternalError creates a new InternalError with the provided details.
// Parameters:
// - location: The location where the error occurred.
// - notification: The notification message for the internal error.
// Returns:
// - An InternalError populated with the given details.
func NewInternalError(location, notification string) InternalError {
	return InternalError{
		Location:     location,
		Notification: notification,
	}
}

// Error implements the error interface for InternalError.
// Returns a formatted string representation of the internal error.
func (err InternalError) Error() string {
	return fmt.Sprintf("location: %s notification: %s", err.Location, err.Notification)
}
