package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

func GetCurrentUserID(ctx *gin.Context) string {
	currentUser := ctx.MustGet("currentUser").(*models.User)
	currentUserID := currentUser.UserID

	return currentUserID
}
