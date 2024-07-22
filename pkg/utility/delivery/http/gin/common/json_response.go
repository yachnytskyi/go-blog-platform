package common

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
)

// GinNewJSONFailureResponse generates a JSON response for a failed operation.
func GinNewJSONFailureResponse(ginContext *gin.Context, err error, httpCode int) {
	jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(err))
	ginContext.JSON(httpCode, jsonResponse)
}

// HandleJSONBindingError handles errors that occur during JSON data binding.
func HandleJSONBindingError(ginContext *gin.Context, logger applicationModel.Logger, location string, err error) {
	httpInternalError := httpError.NewHTTPInternalError(location+".ShouldBindJSON", err.Error())
	logger.Error(httpInternalError)
	GinNewJSONFailureResponse(ginContext, httpInternalError, constants.StatusBadRequest)
}
