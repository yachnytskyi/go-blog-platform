package http

import (
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

type JSONResponseOnSuccess struct {
	Data   any    `json:"data,omitempty"`   // The data field for successful responses.
	Status string `json:"status"`           // The status of the response.
}

type JSONResponseOnFailure struct {
	Error  error  `json:"error,omitempty"`  // The error field for single errors.
	Errors error  `json:"errors,omitempty"` // The errors field for multiple errors.
	Status string `json:"status"`           // The status of the response.
}

func NewJSONResponseOnSuccess(data any) JSONResponseOnSuccess {
	return JSONResponseOnSuccess{
		Data:   data,
		Status: constants.Success,
	}
}

func NewJSONResponseOnFailure(err error) JSONResponseOnFailure {
	jsonResponse := JSONResponseOnFailure{Status: constants.Fail}
	
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
