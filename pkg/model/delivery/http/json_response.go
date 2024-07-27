package http

import (
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	success = "success"
	fail    = "fail"
)

type JSONResponse struct {
	Data   any    `json:"data,omitempty"`   // The data field for successful responses.
	Error  error  `json:"error,omitempty"`  // The error field for single errors.
	Errors error  `json:"errors,omitempty"` // The errors field for multiple errors.
	Status string `json:"status"`           // The status of the response, either "success" or "fail".
}

func NewJSONResponseOnSuccess(data any) JSONResponse {
	return JSONResponse{
		Data:   data,
		Status: success,
	}
}

func NewJSONResponseOnFailure(err error) JSONResponse {
	jsonResponse := JSONResponse{Status: fail}

	switch errorType := err.(type) {
	case model.Errors:
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

func SetStatus(jsonResponse JSONResponse) JSONResponse {
	switch {
	case validator.IsValueNotEmpty(jsonResponse.Data):
		jsonResponse.Status = success
	case validator.IsError(jsonResponse.Error) || validator.IsError(jsonResponse.Errors):
		jsonResponse.Status = fail
	}

	return jsonResponse
}
