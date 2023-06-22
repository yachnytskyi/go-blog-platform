package domain_error

import "fmt"

type DomainValidationError struct {
	Field        string
	FieldType    string
	Notification string
}

func NewDomainValidationError(field string, fieldType string, notification string) error {
	return DomainValidationError{
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

func (err DomainValidationError) Error() string {
	return fmt.Sprintf("field: " + err.Field + " " + "type: " + err.FieldType + " reason: " + err.Notification)
}

type DomainError struct {
	Location     string
	Reason       string
	Notification string
}

func NewDomainError(location string, reason string, notification string) error {
	return DomainError{
		Location:     location,
		Reason:       reason,
		Notification: notification,
	}
}

func (err DomainError) Error() string {
	return fmt.Sprintf("field: " + err.Location + " reason: " + err.Reason + " notification " + err.Notification)
}
