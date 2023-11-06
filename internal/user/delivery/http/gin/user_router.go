package gin

import (
	"github.com/gin-gonic/gin"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type UserRouter struct {
	userController user.UserController
}

func NewUserRouter(userController user.UserController) UserRouter {
	return UserRouter{userController: userController}
}

func (userRouter UserRouter) UserRouter(routerGroup any, userUseCase user.UserUseCase) {
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group("/users")
	router.GET("/", func(ginContext *gin.Context) {
		userRouter.userController.GetAllUsers(ginContext)
	})
	router.GET("/:userID", func(ginContext *gin.Context) {
		userRouter.userController.GetUserById(ginContext)
	})
	router.POST("/login", func(ginContext *gin.Context) {
		userRouter.userController.Login(ginContext)
	})
	router.POST("/register", httpGinMiddleware.AnonymousContextMiddleware(), func(ginContext *gin.Context) {
		userRouter.userController.Register(ginContext)
	})
	router.POST("/forgotten-password", func(ginContext *gin.Context) {
		userRouter.userController.ForgottenPassword(ginContext)
	})
	router.PATCH("/reset-password/:resetToken", func(ginContext *gin.Context) {
		userRouter.userController.ResetUserPassword(ginContext)
	})

	router.Use(httpGinMiddleware.AuthContextMiddleware(userUseCase))
	router.GET("/current_user", func(ginContext *gin.Context) {
		userRouter.userController.GetCurrentUser(ginContext)
	})
	router.PUT("/update", func(ginContext *gin.Context) {
		userRouter.userController.UpdateUserById(ginContext)
	})
	router.DELETE("/delete", func(ginContext *gin.Context) {
		userRouter.userController.Delete(ginContext)
	})
	router.GET("/refresh", func(ginContext *gin.Context) {
		userRouter.userController.RefreshAccessToken(ginContext)
	})
	router.GET("/logout", func(ginContext *gin.Context) {
		userRouter.userController.Logout(ginContext)
	})

}
