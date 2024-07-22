package gin

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type UserRouter struct {
	Logger         applicationModel.Logger
	UserController user.UserController
}

func NewUserRouter(logger applicationModel.Logger, userController user.UserController) UserRouter {
	return UserRouter{
		Logger:         logger,
		UserController: userController,
	}
}

// UserRouter defines the user-related routes and connects them to the corresponding controller methods.
func (userRouter UserRouter) UserRouter(routerGroup any) {
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
	publicAnonymousRoutes.Use(httpGinMiddleware.AnonymousMiddleware(userRouter.Logger))
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
	authenticatedRoutes.Use(httpGinMiddleware.AuthenticationMiddleware(userRouter.Logger))
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
	tokenRoutes.Use(httpGinMiddleware.RefreshTokenAuthenticationMiddleware(userRouter.Logger))
	{
		tokenRoutes.GET(constants.RefreshTokenPath, func(ginContext *gin.Context) {
			userRouter.UserController.RefreshAccessToken(ginContext)
		})

		tokenRoutes.GET(constants.LogoutPath, func(ginContext *gin.Context) {
			userRouter.UserController.Logout(ginContext)
		})
	}
}
