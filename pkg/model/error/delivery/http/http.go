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

type HttpValidationErrorsView []HttpValidationErrorView

func NewHttpValidationErrorsView(httpValidationErrorsView []HttpValidationErrorView) HttpValidationErrorsView {
	return HttpValidationErrorsView(httpValidationErrorsView)
}

func (httpValidationErrorsView HttpValidationErrorsView) Error() string {
	var result strings.Builder
	for _, validationError := range httpValidationErrorsView {
		result.WriteString("field: " + validationError.Field + " " + "type: " + validationError.FieldType + " notification: " + validationError.Notification)
	}
	return result.String()
}

func (httpValidationErrorsView HttpValidationErrorsView) Errors() string {
	var result strings.Builder
	for _, validationError := range httpValidationErrorsView {
		result.WriteString("field: " + validationError.Field + " " + "type: " + validationError.FieldType + " notification: " + validationError.Notification)
	}
	return result.String()
}

type HttpAuthorizationErrorView struct {
	Location     string `json:"-"`
	Notification string `json:"notification"`
}

func NewHttpAuthorizationErrorView(location, notification string) HttpAuthorizationErrorView {
	return HttpAuthorizationErrorView{
		Location:     location,
		Notification: notification,
	}
}

func (err HttpAuthorizationErrorView) Error() string {
	return fmt.Sprintf("location: " + err.Location + " notification: " + err.Notification)
}

type HttpEntityNotFoundErrorView struct {
	Notification string `json:"notification"`
}

func NewHttpEntityNotFoundErrorView(location, notification string) HttpEntityNotFoundErrorView {
	return HttpEntityNotFoundErrorView{
		Notification: notification,
	}
}

func (err HttpEntityNotFoundErrorView) Error() string {
	return fmt.Sprintf("notification: " + err.Notification)
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

type HttpRequestErrorView struct {
	RequestType  string `json:"request_type"`
	Notification string `json:"notification"`
}

func NewHttpRequestErrorView(requestType, notification string) HttpRequestErrorView {
	return HttpRequestErrorView{
		RequestType:  requestType,
		Notification: notification,
	}
}

func (httpRequestErrorView HttpRequestErrorView) Error() string {
	return fmt.Sprintf("request type: " + httpRequestErrorView.RequestType + " notification: " + httpRequestErrorView.Notification)
}

type HttpInternalErrorView struct {
	Location     string `json:"-"`
	Notification string `json:"notification"`
}

func NewHttpInternalErrorView(location, notification string) HttpInternalErrorView {
	return HttpInternalErrorView{
		Notification: notification,
	}
}

func (err HttpInternalErrorView) Error() string {
	return fmt.Sprintf("location " + err.Location + " notification: " + err.Notification)
}
