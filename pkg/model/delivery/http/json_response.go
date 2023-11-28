package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	success = "success"
	fail    = "fail"
)

// JsonResponse represents the structure of an HTTP JSON response.
type JsonResponse struct {
	Data   any    `json:"data,omitempty"`
	Error  error  `json:"error,omitempty"`
	Errors error  `json:"errors,omitempty"`
	Status string `json:"status"`
}

// NewJsonResponseOnSuccess creates a JSON response for a successful operation.
// It sets the "Data" field and the "Status" field to success.
func NewJsonResponseOnSuccess(data any) JsonResponse {
	return JsonResponse{
		Data:   data,
		Status: success,
	}
}

// NewJsonResponseOnFailure creates a JSON response for a failed operation.
// It sets the "Status" field to fail and determines whether to populate "Error" or "Errors" based on the provided error.
func NewJsonResponseOnFailure(err any) JsonResponse {
	jsonResponse := JsonResponse{Status: fail}
	switch errorType := err.(type) {
	case domainError.Errors:
		jsonResponse.Errors = errorType
	case error:
		jsonResponse.Error = errorType
	}
	return jsonResponse
}

// SetStatus sets the "Status" field based on the presence of "Data," "Error," or "Errors."
func SetStatus(jsonResponse *JsonResponse) {
	if jsonResponse.Data != nil {
		jsonResponse.Status = success
	} else if jsonResponse.Error != nil || jsonResponse.Errors != nil {
		jsonResponse.Status = fail
	}
}
