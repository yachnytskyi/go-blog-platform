package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/interfaces"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
)

// HandleJSONBindingError handles errors that occur during JSON data binding.
func HandleJSONBindingError(ginContext *gin.Context, logger interfaces.Logger, location string, err error) {
	httpInternalError := delivery.NewHTTPInternalError(location+".ShouldBindJSON", err.Error())
	logger.Error(httpInternalError)
	ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(httpInternalError)))
}
