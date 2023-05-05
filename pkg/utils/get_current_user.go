package utils

import (
	"github.com/gin-gonic/gin"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
)

func GetCurrentUserID(ctx *gin.Context) string {
	currentUser := ctx.MustGet("currentUser").(*userModel.User)
	currentUserID := currentUser.UserID

	return currentUserID
}
