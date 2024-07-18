package gin

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

type UserRouter struct {
	userController user.UserController
}

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

		publicRoutes.POST(constants.ForgottenPasswordPath, func(ginContext *gin.Context) {
			userRouter.userController.ForgottenPassword(ginContext)
		})

		publicRoutes.PATCH(constants.ResetPasswordPath, func(ginContext *gin.Context) {
			userRouter.userController.ResetUserPassword(ginContext)
		})
	}

	// Public routes with anonymous middleware.
	publicAnonymousRoutes := router.Group("")
	publicAnonymousRoutes.Use(httpGinMiddleware.AnonymousMiddleware())
	{
		publicAnonymousRoutes.POST(constants.LoginPath, func(ginContext *gin.Context) {
			userRouter.userController.Login(ginContext)
		})

		publicAnonymousRoutes.POST(constants.RegisterPath, func(ginContext *gin.Context) {
			userRouter.userController.Register(ginContext)
		})
	}

	// Authenticated routes with authentication middleware.
	authenticatedRoutes := router.Group("")
	authenticatedRoutes.Use(httpGinMiddleware.AuthenticationMiddleware())
	{
		authenticatedRoutes.GET(constants.GetCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.userController.GetCurrentUser(ginContext)
		})

		authenticatedRoutes.PUT(constants.UpdateCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.userController.UpdateCurrentUser(ginContext)
		})

		authenticatedRoutes.DELETE(constants.DeleteCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.userController.DeleteCurrentUser(ginContext)
		})
	}

	// Token-related routes with refresh token authentication middleware.
	tokenRoutes := router.Group("")
	tokenRoutes.Use(httpGinMiddleware.RefreshTokenAuthenticationMiddleware())
	{
		tokenRoutes.GET(constants.RefreshTokenPath, func(ginContext *gin.Context) {
			userRouter.userController.RefreshAccessToken(ginContext)
		})

		tokenRoutes.GET(constants.LogoutPath, func(ginContext *gin.Context) {
			userRouter.userController.Logout(ginContext)
		})
	}
}
