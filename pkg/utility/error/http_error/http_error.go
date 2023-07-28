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

type ErrorMessageView struct {
	Notification string
}

func NewErrorMessage(notification string) error {
	return ErrorMessageView{
		Notification: notification,
	}
}

func (err ErrorMessageView) Error() string {
	return fmt.Sprintf("notification: " + err.Notification)
}
