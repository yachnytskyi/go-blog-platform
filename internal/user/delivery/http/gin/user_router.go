package gin

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

// Constants for route paths.
const (
	registerPath          = "/register"
	forgottenPasswordPath = "/forgotten-password"
	resetPasswordPath     = "/reset-password/:resetToken"
	loginPath             = "/login"
	getCurrentUserPath    = "/current_user"
	updateCurrentUserPath = "/update"
	deleteCurrentUserPath = "/delete"
	refreshTokenPath      = "/refresh"
	logoutPath            = "/logout"
)

// UserRouter is responsible for defining the user-related routes and handling HTTP requests.
type UserRouter struct {
	userController user.UserController
}

// NewUserRouter creates a new instance of UserRouter with the provided user controller.
func NewUserRouter(userController user.UserController) UserRouter {
	return UserRouter{userController: userController}
}

// UserRouter defines the user-related routes and connects them to the corresponding controller methods.
func (userRouter UserRouter) UserRouter(routerGroup any) {
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group(constants.UsersGroupPath)

	// Public routes.
	publicRoutes := router.Group("")
	{
		publicRoutes.GET(constants.GetAllItemsURL, func(ginContext *gin.Context) {
			userRouter.userController.GetAllUsers(ginContext)
		})

		publicRoutes.GET(constants.GetItemByIdURL, func(ginContext *gin.Context) {
			userRouter.userController.GetUserById(ginContext)
		})

		publicRoutes.POST(forgottenPasswordPath, func(ginContext *gin.Context) {
			userRouter.userController.ForgottenPassword(ginContext)
		})

		publicRoutes.PATCH((resetPasswordPath), func(ginContext *gin.Context) {
			userRouter.userController.ResetUserPassword(ginContext)
		})
	}

	// Public routes with anonymous middleware.
	publicAnonymousRoutes := router.Group("")
	publicAnonymousRoutes.Use(httpGinMiddleware.AnonymousMiddleware())
	{
		publicAnonymousRoutes.POST(loginPath, func(ginContext *gin.Context) {
			userRouter.userController.Login(ginContext)
		})

		publicAnonymousRoutes.POST(registerPath, func(ginContext *gin.Context) {
			userRouter.userController.Register(ginContext)
		})

	}

	// Authenticated routes with authentication middleware.
	authenticatedRoutes := router.Group("")
	authenticatedRoutes.Use(httpGinMiddleware.AuthenticationMiddleware())
	{
		authenticatedRoutes.GET(getCurrentUserPath, httpGinMiddleware.AuthenticationMiddleware(), func(ginContext *gin.Context) {
			userRouter.userController.GetCurrentUser(ginContext)
		})

		authenticatedRoutes.PUT(updateCurrentUserPath, httpGinMiddleware.AuthenticationMiddleware(), func(ginContext *gin.Context) {
			userRouter.userController.UpdateCurrentUser(ginContext)
		})

		authenticatedRoutes.DELETE(deleteCurrentUserPath, httpGinMiddleware.AuthenticationMiddleware(), func(ginContext *gin.Context) {
			userRouter.userController.DeleteCurrentUser(ginContext)
		})
	}

	// Token-related routes with refresh token authentication middleware.
	tokenRoutes := router.Group("")
	tokenRoutes.Use(httpGinMiddleware.RefreshTokenAuthenticationMiddleware())
	{
		tokenRoutes.GET(refreshTokenPath, httpGinMiddleware.RefreshTokenAuthenticationMiddleware(), func(ginContext *gin.Context) {
			userRouter.userController.RefreshAccessToken(ginContext)
		})

		tokenRoutes.GET(logoutPath, httpGinMiddleware.RefreshTokenAuthenticationMiddleware(), func(ginContext *gin.Context) {
			userRouter.userController.Logout(ginContext)
		})
	}
}
