package http

import (
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	success = "success" // Constant representing a successful operation.
	fail    = "fail"    // Constant representing a failed operation.
)

// JSONResponse represents the structure of an HTTP JSON response.
type JSONResponse struct {
	Data   any    `json:"data,omitempty"`   // The data field for successful responses.
	Error  error  `json:"error,omitempty"`  // The error field for single errors.
	Errors error  `json:"errors,omitempty"` // The errors field for multiple errors.
	Status string `json:"status"`           // The status of the response, either "success" or "fail".
}

// NewJSONSuccessResponse creates a JSON response for a successful operation.
// Parameters:
// - data: The data to be included in the response.
// Returns:
// - A JSONResponse with the "Data" field populated and the "Status" set to "success".
func NewJSONSuccessResponse(data any) JSONResponse {
	return JSONResponse{
		Data:   data,
		Status: success,
	}
}

// NewJSONFailureResponse creates a JSON response for a failed operation.
// Parameters:
// - err: The error to be included in the response.
// Returns:
// - A JSONResponse with the "Status" set to "fail" and the "Error" or "Errors" field populated based on the type of error.
func NewJSONFailureResponse(err error) JSONResponse {
	jsonResponse := JSONResponse{Status: fail}

	switch errorType := err.(type) {
	case domainModel.Errors:
		if errorType.Len() == 1 {
			jsonResponse.Error = errorType
		} else {
			jsonResponse.Errors = errorType
		}
	case error:
		jsonResponse.Error = errorType
	default:
		jsonResponse.Error = errorType
	}

	return jsonResponse
}

// SetStatus sets the "Status" field based on the presence of "Data," "Error," or "Errors."
// Parameters:
// - jsonResponse: A pointer to the JSONResponse whose status needs to be set.
func SetStatus(jsonResponse *JSONResponse) {
	switch {
	case validator.IsValueNotEmpty(jsonResponse.Data):
		jsonResponse.Status = success
	case validator.IsError(jsonResponse.Error) || validator.IsError(jsonResponse.Errors):
		// Optionally handle cases where neither Data nor Error/Errors are set.
		jsonResponse.Status = fail
	}
}
