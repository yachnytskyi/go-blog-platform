package gin

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

// Constants for route paths.
const (
	groupPath             = "/users"
	getAllUsersPath       = ""
	getUserByIdPath       = "/:"
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
func (userRouter UserRouter) UserRouter(routerGroup any, userUseCase user.UserUseCase) {
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group(groupPath)

	// Public routes.
	router.GET(constants.EmptyString, func(ginContext *gin.Context) {
		userRouter.userController.GetAllUsers(ginContext)
	})

	router.GET(getUserByIdPath+constants.UserIDContext, func(ginContext *gin.Context) {
		userRouter.userController.GetUserById(ginContext)
	})

	router.POST(registerPath, httpGinMiddleware.AnonymousMiddleware(), func(ginContext *gin.Context) {
		userRouter.userController.Register(ginContext)
	})

	router.POST(loginPath, func(ginContext *gin.Context) {
		userRouter.userController.Login(ginContext)
	})

	router.POST(forgottenPasswordPath, func(ginContext *gin.Context) {
		userRouter.userController.ForgottenPassword(ginContext)
	})

	router.PATCH((resetPasswordPath), func(ginContext *gin.Context) {
		userRouter.userController.ResetUserPassword(ginContext)
	})

	// Authenticated routes.
	router.GET(getCurrentUserPath, httpGinMiddleware.AuthenticationMiddleware(userUseCase), func(ginContext *gin.Context) {
		userRouter.userController.GetCurrentUser(ginContext)
	})

	router.PUT(updateCurrentUserPath, httpGinMiddleware.AuthenticationMiddleware(userUseCase), func(ginContext *gin.Context) {
		userRouter.userController.UpdateCurrentUser(ginContext)
	})

	router.DELETE(deleteCurrentUserPath, httpGinMiddleware.AuthenticationMiddleware(userUseCase), func(ginContext *gin.Context) {
		userRouter.userController.DeleteCurrentUser(ginContext)
	})

	// Token-related routes.
	router.GET(refreshTokenPath, httpGinMiddleware.RefreshTokenAuthenticationMiddleware(userUseCase), func(ginContext *gin.Context) {
		userRouter.userController.RefreshAccessToken(ginContext)
	})

	router.GET(logoutPath, httpGinMiddleware.RefreshTokenAuthenticationMiddleware(userUseCase), func(ginContext *gin.Context) {
		userRouter.userController.Logout(ginContext)
	})
}
