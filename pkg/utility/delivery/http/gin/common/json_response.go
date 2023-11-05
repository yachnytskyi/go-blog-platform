package common

import (
	"github.com/gin-gonic/gin"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
)

// GinNewJsonResponseOnFailure generates a JSON response for a failed operation.
// It takes a Gin context, an error that occurred, and an HTTP status code.
func GinNewJsonResponseOnFailure(ginContext *gin.Context, err error, httpCode int) {
	jsonResponse := httpModel.NewJsonResponseOnFailure(httpError.HandleError(err))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(httpCode, jsonResponse)
}
