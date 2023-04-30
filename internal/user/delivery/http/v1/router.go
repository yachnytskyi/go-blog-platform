package v1

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

func (userRouter *UserRouter) UserRouter(routerGroup *gin.RouterGroup, userService user.Service) {
	router := routerGroup.Group("/users")

	router.POST("/register", userRouter.userHandler.Register)
	router.POST("/login", userRouter.userHandler.Login)

	router.Use(middleware.DeserializeUser(userService))
	router.POST("/forgotten-password", userRouter.userHandler.ForgottenPassword)
	router.PATCH("/reset-password/:resetToken", userRouter.userHandler.ResetUserPassword)

	router.GET("/refresh", userRouter.userHandler.RefreshAccessToken)
	router.GET("/logout", userRouter.userHandler.Logout)

	router.GET("/me", userRouter.userHandler.GetMe)
	router.PUT("/update", userRouter.userHandler.UpdateUserById)
	router.DELETE("/delete", userRouter.userHandler.Delete)
}
