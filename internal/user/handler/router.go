package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/middleware"
)

type UserRouter struct {
	userHandler UserHandler
}

func NewUserRouter(userHandler UserHandler) UserRouter {
	return UserRouter{userHandler: userHandler}
}

func (userRouter *UserRouter) UserRoute(rg *gin.RouterGroup, userService user.Service) {
	router := rg.Group("/users")

	router.POST("/register", userRouter.userHandler.Register)
	router.POST("/login", userRouter.userHandler.Login)
	router.GET("/refresh", userRouter.userHandler.RefreshAccessToken)
	router.GET("/logout", middleware.DeserializeUser(userService), userRouter.userHandler.LogoutUser)

	router.Use(middleware.DeserializeUser(userService))
	router.GET("/me", userRouter.userHandler.GetMe)
}
