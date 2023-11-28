package domain

import (
	"fmt"
	"strings"
)

type Errors interface {
	error
	Errors() string
}

type ValidationError struct {
	Location     string
	Field        string
	FieldType    string
	Notification string
}

func NewValidationError(location, field, fieldType, notification string) ValidationError {
	return ValidationError{
		Location:     location,
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("location: " + err.Location + " field: " + err.Field + " " + "type: " + err.FieldType + " notification: " + err.Notification)
}

type ValidationErrors []ValidationError

func NewValidationErrors(validationErrors []ValidationError) ValidationErrors {
	return ValidationErrors(validationErrors)
}

func (validationErrors ValidationErrors) Error() string {
	var result strings.Builder
	for _, validationError := range validationErrors {
		result.WriteString("field: " + validationError.Field + " " + "type: " + validationError.FieldType + " notification: " + validationError.Notification)

	}
	return result.String()
}

func (validationErrors ValidationErrors) Errors() string {
	var result strings.Builder
	for _, validationError := range validationErrors {
		result.WriteString("field: " + validationError.Field + " " + "type: " + validationError.FieldType + " notification: " + validationError.Notification)

	}
	return result.String()
}

type AuthorizationError struct {
	Location     string
	Notification string
}

func NewAuthorizationError(location, notification string) AuthorizationError {
	return AuthorizationError{
		Location:     location,
		Notification: notification,
	}
}

func (err AuthorizationError) Error() string {
	return fmt.Sprintf("location: " + err.Location + " notification: " + err.Notification)
}

type EntityNotFoundError struct {
	Location     string
	Query        string
	Notification string
}

func NewEntityNotFoundError(location, query, notification string) EntityNotFoundError {
	return EntityNotFoundError{
		Location:     location,
		Query:        query,
		Notification: notification,
	}
}

func (err EntityNotFoundError) Error() string {
	return fmt.Sprintf("location: " + err.Location + " query: " + err.Query + " notification: " + err.Notification)
}

type PaginationError struct {
	CurrentPage  string
	TotalPages   string
	Notification string
}

func NewPaginationError(currentPage, totalPages, notification string) PaginationError {
	return PaginationError{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		Notification: notification,
	}
}

func (paginationError PaginationError) Error() string {
	return fmt.Sprintf("page: " + paginationError.CurrentPage + " total: " +
		paginationError.TotalPages + " notification: " + paginationError.Notification)
}

type InternalError struct {
	Location     string
	Notification string
}

func NewInternalError(location, notification string) InternalError {
	return InternalError{
		Location:     location,
		Notification: notification,
	}
}

func (err InternalError) Error() string {
	return fmt.Sprintf("location: " + err.Location + " notification: " + err.Notification)
}
