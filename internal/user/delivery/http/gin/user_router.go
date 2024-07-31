package gin

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type UserRouter struct {
	interfaces.Config
	Logger         interfaces.Logger
	UserController interfaces.UserController
}

func NewUserRouter(config interfaces.Config, logger interfaces.Logger, userController interfaces.UserController) UserRouter {
	return UserRouter{
		Config:         config,
		Logger:         logger,
		UserController: userController,
	}
}

// UserRouter defines the user-related routes and connects them to the corresponding controller methods.
func (userRouter UserRouter) Router(routerGroup any) {
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group(constants.UsersGroupPath)

	// Public routes.
	publicRoutes := router.Group("")
	{
		publicRoutes.GET(constants.GetAllItemsURL, func(ginContext *gin.Context) {
			userRouter.UserController.GetAllUsers(ginContext)
		})

		publicRoutes.GET(constants.GetItemByIdURL, func(ginContext *gin.Context) {
			userRouter.UserController.GetUserById(ginContext)
		})

		publicRoutes.POST(constants.ForgottenPasswordPath, func(ginContext *gin.Context) {
			userRouter.UserController.ForgottenPassword(ginContext)
		})

		publicRoutes.PATCH(constants.ResetPasswordPath, func(ginContext *gin.Context) {
			userRouter.UserController.ResetUserPassword(ginContext)
		})
	}

	// Public routes with anonymous middleware.
	publicAnonymousRoutes := router.Group("")
	publicAnonymousRoutes.Use(middleware.AnonymousMiddleware(userRouter.Logger))
	{
		publicAnonymousRoutes.POST(constants.LoginPath, func(ginContext *gin.Context) {
			userRouter.UserController.Login(ginContext)
		})

		publicAnonymousRoutes.POST(constants.RegisterPath, func(ginContext *gin.Context) {
			userRouter.UserController.Register(ginContext)
		})
	}

	// Authenticated routes with authentication middleware.
	authenticatedRoutes := router.Group("")
	authenticatedRoutes.Use(middleware.AuthenticationMiddleware(userRouter.Config, userRouter.Logger))
	{
		authenticatedRoutes.GET(constants.GetCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.UserController.GetCurrentUser(ginContext)
		})

		authenticatedRoutes.PUT(constants.UpdateCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.UserController.UpdateCurrentUser(ginContext)
		})

		authenticatedRoutes.DELETE(constants.DeleteCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.UserController.DeleteCurrentUser(ginContext)
		})
	}

	// Token-related routes with refresh token authentication middleware.
	tokenRoutes := router.Group("")
	tokenRoutes.Use(middleware.RefreshTokenAuthenticationMiddleware(userRouter.Config, userRouter.Logger))
	{
		tokenRoutes.GET(constants.RefreshTokenPath, func(ginContext *gin.Context) {
			userRouter.UserController.RefreshAccessToken(ginContext)
		})

		tokenRoutes.GET(constants.LogoutPath, func(ginContext *gin.Context) {
			userRouter.UserController.Logout(ginContext)
		})
	}
}
