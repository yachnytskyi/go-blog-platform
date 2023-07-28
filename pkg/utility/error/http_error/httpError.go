package http_error

import "fmt"

type HttpValidationErrorView struct {
	Field        string `json:"field"`
	FieldType    string `json:"type"`
	Notification string `json:"notification"`
}

func NewHttpValidationError(field string, fieldType string, notification string) error {
	return HttpValidationErrorView{
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

func (err HttpValidationErrorView) Error() string {
	return fmt.Sprintf("field: " + err.Field + " " + "type: " + err.FieldType + " notification: " + err.Notification)
}

type EntityNotFoundErrorView struct {
	Location string
	Code     string
	Reason   string
}

func NewEntityNotFoundError(location string, reason string) error {
	return EntityNotFoundErrorView{
		Location: location,
		Reason:   reason,
	}
}

func (err EntityNotFoundErrorView) Error() string {
	return fmt.Sprintf("field: " + err.Location + " reason: " + err.Reason)
}

type ErrorMessage struct {
	Notification string
}

func NewErrorMessage(notification string) error {
	return ErrorMessage{
		Notification: notification,
	}
}

func (err ErrorMessage) Error() string {
	return fmt.Sprintf("notification: " + err.Notification)
}
