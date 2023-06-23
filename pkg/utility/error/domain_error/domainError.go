package domain_error

import "fmt"

type ValidationError struct {
	Field        string
	FieldType    string
	Notification string
}

func NewValidationError(field string, fieldType string, notification string) error {
	return ValidationError{
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

func (err ValidationError) Error() string {
	return fmt.Sprintf("field: " + err.Field + " " + "type: " + err.FieldType + " reason: " + err.Notification)
}

type InternalError struct {
	Location     string
	Reason       string
	Notification string
}

func NewInternalError(location string, reason string, notification string) error {
	return InternalError{
		Location:     location,
		Reason:       reason,
		Notification: notification,
	}
}

func (err InternalError) Error() string {
	return fmt.Sprintf("field: " + err.Location + " reason: " + err.Reason + " notification " + err.Notification)
}
