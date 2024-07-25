package common

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
)

// HandleJSONBindingError handles errors that occur during JSON data binding.
func HandleJSONBindingError(ginContext *gin.Context, logger model.Logger, location string, err error) {
	httpInternalError := httpError.NewHTTPInternalError(location+".ShouldBindJSON", err.Error())
	logger.Error(httpInternalError)
	ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(err)))
}
