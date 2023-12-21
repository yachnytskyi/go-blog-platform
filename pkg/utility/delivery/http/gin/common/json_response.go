package common

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

// GinNewJSONResponseOnFailure generates a JSON response for a failed operation.
// It takes a Gin context, an error that occurred, and an HTTP status code.
func GinNewJSONResponseOnFailure(ginContext *gin.Context, err error, httpCode int) {
	jsonResponse := httpModel.NewJSONResponseOnFailure(httpError.HandleError(err))
	ginContext.JSON(httpCode, jsonResponse)
}

// HandleJSONBindingError handles errors that occur during JSON data binding.
// It logs the error, generates a JSON response, and sets the HTTP status code.
func HandleJSONBindingError(ginContext *gin.Context, location string, err error) {
	internalError := httpError.NewHttpInternalErrorView(location+".ShouldBindJSON", err.Error())
	logging.Logger(internalError)
	GinNewJSONResponseOnFailure(ginContext, internalError, constants.StatusBadRequest)
}
