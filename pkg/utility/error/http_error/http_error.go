package http_error

import (
	"fmt"
	"strings"
)

type HttpValidationBaseErrorView struct {
	Field        string `json:"field"`
	FieldType    string `json:"type"`
	Notification string `json:"notification"`
	Status       string `json:"status,omitempty"`
}

func NewHttpValidationError(field string, fieldType string, notification string) error {
	return &HttpValidationBaseErrorView{
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

func (httpValidationBaseErrorView *HttpValidationBaseErrorView) Error() string {
	return fmt.Sprintf("field: " + httpValidationBaseErrorView.Field + " " + "type: " +
		httpValidationBaseErrorView.FieldType + " notification: " + httpValidationBaseErrorView.Notification + httpValidationBaseErrorView.Status)
}

type HttpValidationErrorView struct {
	HttpValidationErrorView *HttpValidationBaseErrorView `json:"error"`
}

func NewHttpValidationBaseError(httpValidationBaseErrorView *HttpValidationBaseErrorView) error {
	return &HttpValidationErrorView{
		HttpValidationErrorView: httpValidationBaseErrorView,
	}
}

func (httpValidationErrorView *HttpValidationErrorView) Error() string {
	return fmt.Sprintf("field: " + httpValidationErrorView.HttpValidationErrorView.Field + " " + "type: " +
		httpValidationErrorView.HttpValidationErrorView.FieldType + " notification: " + httpValidationErrorView.HttpValidationErrorView.Notification)
}

type HttpValidationErrorsView struct {
	HttpValidationErrorsView []*HttpValidationBaseErrorView `json:"errors"`
	Status                   string                         `json:"status"`
}

func NewHttpValidationErrorsView(httpValidationErrorsView []*HttpValidationBaseErrorView) error {
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

type HttpBaseErrorMessageView struct {
	Notification string `json:"error"`
	Status       string `json:"status,omitempty"`
}

func NewHttpBaseErrorMessage(notification string) error {
	return &HttpBaseErrorMessageView{
		Notification: notification,
	}
}

func (httpMessageErrorView *HttpBaseErrorMessageView) Error() string {
	return fmt.Sprintf("notification: " + httpMessageErrorView.Notification + httpMessageErrorView.Status)
}

type HttpErrorMessageView struct {
	HttpErrorMessageView *HttpBaseErrorMessageView `json:"error"`
}

func NewHttpErrorMessage(httpBaseErrorMessageView *HttpBaseErrorMessageView) error {
	return &HttpErrorMessageView{
		HttpErrorMessageView: httpBaseErrorMessageView,
	}
}

func (httpMessageErrorView *HttpErrorMessageView) Error() string {
	return fmt.Sprintf("notification: " + httpMessageErrorView.HttpErrorMessageView.Notification)
}
