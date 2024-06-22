package common

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

// GinNewJSONFailureResponse generates a JSON response for a failed operation.
// It takes a Gin context, an error that occurred, and an HTTP status code.
// Parameters:
// - ginContext: The Gin context used to generate the JSON response.
// - err: The error that occurred.
// - httpCode: The HTTP status code to be set in the response.
func GinNewJSONFailureResponse(ginContext *gin.Context, err error, httpCode int) {
	jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(err))
	ginContext.JSON(httpCode, jsonResponse)
}

// HandleJSONBindingError handles errors that occur during JSON data binding.
// It logs the error, generates a JSON response, and sets the HTTP status code.
// Parameters:
// - ginContext: The Gin context used to generate the JSON response.
// - location: A string representing the location or context for error logging.
// - err: The error that occurred during JSON binding.
func HandleJSONBindingError(ginContext *gin.Context, location string, err error) {
	httpInternalError := httpError.NewHTTPInternalError(location+".ShouldBindJSON", err.Error())
	logging.Logger(httpInternalError)
	GinNewJSONFailureResponse(ginContext, httpInternalError, constants.StatusBadRequest)
}
