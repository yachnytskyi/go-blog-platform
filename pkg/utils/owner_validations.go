package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

func GetCurrentUserID(ctx *gin.Context) string {
	currentUser := ctx.MustGet("currentUser").(*models.UserFullResponse)
	currentUserID := currentUser.UserID

	return currentUserID
}

func IsUserOwner(currentUserID string, userID string) error {
	if currentUserID != userID {
		return errors.New("sorry, but you do not have permissions to do that")
	}

	return nil
}
