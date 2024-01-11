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

// NewJSONSuccessResponse creates a JSON response for a successful JSON operation.
// It sets the "Data" field and the "Status" field to success.
func NewJSONSuccessResponse(data any) JSONResponse {
	return JSONResponse{
		Data:   data,
		Status: success,
	}
}

// NewJSONFailureResponse creates a JSON response for a failed JSON operation.
// It sets the "Status" field to fail and determines whether to populate "Error" or "Errors" based on the provided error.
func NewJSONFailureResponse(err error) JSONResponse {
	jsonResponse := JSONResponse{Status: fail}
	switch errorType := err.(type) {
	case domainError.Errors:
		if errorType.Len(errorType) == 1 {
			jsonResponse.Error = errorType
		} else {
			jsonResponse.Errors = errorType
		}
	case error:
		jsonResponse.Error = errorType
	}
	return jsonResponse
}

// SetStatus sets the "Status" field based on the presence of "Data," "Error," or "Errors."
func SetStatus(jsonResponse *JSONResponse) {
	if validator.IsValueNotEmpty(jsonResponse.Data) {
		jsonResponse.Status = success
	} else if validator.IsError(jsonResponse.Error) || validator.IsError(jsonResponse.Errors) {
		jsonResponse.Status = fail
	}
}
