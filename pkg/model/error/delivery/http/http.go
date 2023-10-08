package http

import (
	"fmt"
	"strings"
)

type HttpValidationErrorView struct {
	Field        string `json:"field"`
	FieldType    string `json:"type"`
	Notification string `json:"notification"`
}

func NewHttpValidationError(field string, fieldType string, notification string) HttpValidationErrorView {
	return HttpValidationErrorView{
		Field:        field,
		FieldType:    fieldType,
		Notification: notification,
	}
}

func (httpValidationErrorView HttpValidationErrorView) Error() string {
	return fmt.Sprintf("field: " + httpValidationErrorView.Field + " " + "type: " +
		httpValidationErrorView.FieldType + " notification: " + httpValidationErrorView.Notification)
}

type HttpValidationErrorsView struct {
	HttpValidationErrorsView []HttpValidationErrorView `json:"errors"`
}

func NewHttpValidationErrorsView(httpValidationErrorsView []HttpValidationErrorView) HttpValidationErrorsView {
	return HttpValidationErrorsView{
		HttpValidationErrorsView: httpValidationErrorsView,
	}
}

func (httpValidationErrorsView HttpValidationErrorsView) Error() string {
	var result strings.Builder
	for _, validationError := range httpValidationErrorsView.HttpValidationErrorsView {
		result.WriteString("field: " + validationError.Field + " " + "type: " + validationError.FieldType + " notification: " + validationError.Notification)
	}
	return result.String()
}

type HttpErrorMessageView struct {
	Notification string `json:"notification"`
}

func NewHttpErrorMessage(notification string) HttpErrorMessageView {
	return HttpErrorMessageView{
		Notification: notification,
	}
}

func (httpMessageErrorView HttpErrorMessageView) Error() string {
	return fmt.Sprintf("notification: " + httpMessageErrorView.Notification)
}

type HttpPaginationErrorView struct {
	CurrentPage  string `json:"current_page"`
	TotalPages   string `json:"total_pages"`
	Notification string `json:"notification"`
}

func NewHttpPaginationErrorView(currentPage, totalPages, notification string) HttpPaginationErrorView {
	return HttpPaginationErrorView{
		CurrentPage:  currentPage,
		TotalPages:   totalPages,
		Notification: notification,
	}
}

func (httpPaginationErrorView HttpPaginationErrorView) Error() string {
	return fmt.Sprintf("current page: " + httpPaginationErrorView.CurrentPage + " total pages: " +
		httpPaginationErrorView.TotalPages + " notification: " + httpPaginationErrorView.Notification)
}
