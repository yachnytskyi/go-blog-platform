package http_error

import (
	"fmt"
	"strings"
)

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

type HttpValidationErrorsView struct {
	HttpValidationErrorsView []*HttpValidationErrorView `json:"errors"`
}

func NewHttpValidationErrorsView(httpValidationErrorsView []*HttpValidationErrorView) error {
	return &HttpValidationErrorsView{
		HttpValidationErrorsView: httpValidationErrorsView,
	}
}

func (httpValidationErrorsView *HttpValidationErrorsView) Error() string {
	var result strings.Builder
	for _, vavalidationError := range httpValidationErrorsView.HttpValidationErrorsView {
		result.WriteString("field: " + vavalidationError.Field + " " + "type: " + vavalidationError.FieldType + " notification: " + vavalidationError.Notification)

	}

	return result.String()
}

type HttpErrorMessageView struct {
	Notification string
}

func NewErrorMessage(notification string) error {
	return HttpErrorMessageView{
		Notification: notification,
	}
}

func (err HttpErrorMessageView) Error() string {
	return fmt.Sprintf("notification: " + err.Notification)
}
