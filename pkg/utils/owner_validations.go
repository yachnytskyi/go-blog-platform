package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

func IsOwner(userID string, ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.UserDB)
	currentUserID := currentUser.UserID.Hex()
	println(userID)
	println(currentUserID)
	println(currentUser.UserID.Hex())

	if currentUserID != userID {
		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
		println("Error")
		return
	}
}
