package http

import (
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	success = "success"
	fail    = "fail"
)

// JSONResponse represents the structure of an HTTP JSON response.
type JSONResponse struct {
	Data   any    `json:"data,omitempty"`
	Error  error  `json:"error,omitempty"`
	Errors error  `json:"errors,omitempty"`
	Status string `json:"status"`
}

// NewJSONResponseOnSuccess creates a JSON response for a successful operation.
// It sets the "Data" field and the "Status" field to success.
func NewJSONResponseOnSuccess(data any) JSONResponse {
	return JSONResponse{
		Data:   data,
		Status: success,
	}
}

// NewJSONResponseOnFailure creates a JSON response for a failed operation.
// It sets the "Status" field to fail and determines whether to populate "Error" or "Errors" based on the provided error.
func NewJSONResponseOnFailure(err error) JSONResponse {
	jsonResponse := JSONResponse{Status: fail}
	switch errorType := err.(type) {
	case domainError.Errors:
		jsonResponse.Errors = errorType
	case error:
		jsonResponse.Error = errorType
	}
	return jsonResponse
}

// SetStatus sets the "Status" field based on the presence of "Data," "Error," or "Errors."
func SetStatus(jsonResponse *JSONResponse) {
	if validator.IsValueNotNil(jsonResponse.Data) {
		jsonResponse.Status = success
	} else if validator.IsErrorNotNil(jsonResponse.Error) || validator.IsErrorNotNil(jsonResponse.Errors) {
		jsonResponse.Status = fail
	}
}
