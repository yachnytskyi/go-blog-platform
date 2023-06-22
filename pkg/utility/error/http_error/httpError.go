package http_error

import "fmt"

type HttpValidationError struct {
	Field        string `json:"field"`
	FieldType    string `json:"type"`
	Notification string `json:"notification"`
}

func NewHttpValidationError(field string, fieldType string, notification string) error {
	return HttpValidationError{
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

func (err HttpValidationError) Error() string {
	return fmt.Sprintf("field: " + err.Field + " " + "type: " + err.FieldType + " notification: " + err.Notification)
}

type HttpError struct {
	Notification string `json:"notification"`
}

func NewHttpError(notification string) error {
	return HttpValidationError{
		Notification: notification,
	}
}

func (err HttpError) Error() string {
	return fmt.Sprintf("notification: " + err.Notification)
}
