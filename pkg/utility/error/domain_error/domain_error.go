package domain_error

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field        string
	FieldType    string
	Notification string
}

func NewValidationError(field string, fieldType string, notification string) error {
	return &ValidationError{
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

func (err *ValidationError) Error() string {
	return fmt.Sprintf("field: " + err.Field + " " + "type: " + err.FieldType + " notification: " + err.Notification)
}

type ValidationErrors struct {
	ValidationErrors []*ValidationError
}

func NewValidationErrors(validationErrors []*ValidationError) error {
	return &ValidationErrors{
		ValidationErrors: validationErrors,
	}
}

func (validationErrors *ValidationErrors) Error() string {
	var result strings.Builder
	for _, vavalidationError := range validationErrors.ValidationErrors {
		result.WriteString("field: " + vavalidationError.Field + " " + "type: " + vavalidationError.FieldType + " notification: " + vavalidationError.Notification)

	}

	return result.String()
}

type EntityNotFoundError struct {
	Location string
	Code     string
	Reason   string
}

func NewEntityNotFoundError(location string, reason string) error {
	return &EntityNotFoundError{
		Location: location,
		Reason:   reason,
	}
}

func (err *EntityNotFoundError) Error() string {
	return fmt.Sprintf("field: " + err.Location + " reason: " + err.Reason)
}

type InternalError struct {
	Location string
	Code     string
	Reason   string
}

func NewInternalError(location string, reason string) error {
	return &InternalError{
		Location: location,
		Reason:   reason,
	}
}

func (err *InternalError) Error() string {
	return fmt.Sprintf("field: " + err.Location + " reason: " + err.Reason)
}

type ErrorMessage struct {
	Notification string
}

func NewErrorMessage(notification string) error {
	return &ErrorMessage{
		Notification: notification,
	}
}

func (err *ErrorMessage) Error() string {
	return fmt.Sprintf("notification: " + err.Notification)
}
