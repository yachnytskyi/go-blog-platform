package common

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

// GinNewJSONFailureResponse generates a JSON response for a failed operation.
func GinNewJSONFailureResponse(ginContext *gin.Context, err error, httpCode int) {
	jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(err))
	ginContext.JSON(httpCode, jsonResponse)
}

// HandleJSONBindingError handles errors that occur during JSON data binding.
func HandleJSONBindingError(ginContext *gin.Context, location string, err error) {
	httpInternalError := httpError.NewHTTPInternalError(location+".ShouldBindJSON", err.Error())
	logging.Logger(httpInternalError)
	GinNewJSONFailureResponse(ginContext, httpInternalError, constants.StatusBadRequest)
}
