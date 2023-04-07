package user

import "github.com/gin-gonic/gin"

type Handler interface {
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	UpdateUser(ctx *gin.Context)
	RefreshAccessToken(ctx *gin.Context)
	GetMe(ctx *gin.Context)
	Logout(ctx *gin.Context)
	Delete(ctx *gin.Context)
}
