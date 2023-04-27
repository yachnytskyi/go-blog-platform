package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

func IsOwner(ctx *gin.Context, userID string) error {
	currentUser := ctx.MustGet("currentUser").(*models.UserDB)
	currentUserID := currentUser.UserID.Hex()

	if currentUserID != userID {
		return errors.New("sorry, but you cannot edit this post")
	}

	return nil
}
