package gin

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
)

// UserRouter is responsible for defining the user-related routes and handling HTTP requests.
// It includes routes for public access, authentication, and token management.
type UserRouter struct {
	userController user.UserController
}

// NewUserRouter creates a new instance of UserRouter with the provided user controller.
//
// Parameters:
// - userController (user.UserController): The controller to handle user-related requests.
//
// Returns:
// - UserRouter: A new UserRouter instance.
func NewUserRouter(userController user.UserController) UserRouter {
	return UserRouter{userController: userController}
}

// UserRouter defines the user-related routes and connects them to the corresponding controller methods.
// It sets up routes for public access, routes that require authentication, and routes for token management.
//
// Parameters:
// - routerGroup (any): The router group to which the user routes will be added.
func (userRouter UserRouter) UserRouter(routerGroup any) {
	// Type assertion to map the generic routerGroup to a Gin RouterGroup.
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group(constants.UsersGroupPath)

	// Public routes.
	publicRoutes := router.Group("")
	{
		// GET all users.
		publicRoutes.GET(constants.GetAllItemsURL, func(ginContext *gin.Context) {
			userRouter.userController.GetAllUsers(ginContext)
		})

		// GET user by ID.
		publicRoutes.GET(constants.GetItemByIdURL, func(ginContext *gin.Context) {
			userRouter.userController.GetUserById(ginContext)
		})

		// POST request for forgotten password.
		publicRoutes.POST(constants.ForgottenPasswordPath, func(ginContext *gin.Context) {
			userRouter.userController.ForgottenPassword(ginContext)
		})

		// PATCH request for resetting password.
		publicRoutes.PATCH(constants.ResetPasswordPath, func(ginContext *gin.Context) {
			userRouter.userController.ResetUserPassword(ginContext)
		})
	}

	// Public routes with anonymous middleware.
	publicAnonymousRoutes := router.Group("")
	publicAnonymousRoutes.Use(httpGinMiddleware.AnonymousMiddleware())
	{
		// POST request for user login.
		publicAnonymousRoutes.POST(constants.LoginPath, func(ginContext *gin.Context) {
			userRouter.userController.Login(ginContext)
		})

		// POST request for user registration.
		publicAnonymousRoutes.POST(constants.RegisterPath, func(ginContext *gin.Context) {
			userRouter.userController.Register(ginContext)
		})
	}

	// Authenticated routes with authentication middleware.
	authenticatedRoutes := router.Group("")
	authenticatedRoutes.Use(httpGinMiddleware.AuthenticationMiddleware())
	{
		// GET request to fetch current user details.
		authenticatedRoutes.GET(constants.GetCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.userController.GetCurrentUser(ginContext)
		})

		// PUT request to update current user details.
		authenticatedRoutes.PUT(constants.UpdateCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.userController.UpdateCurrentUser(ginContext)
		})

		// DELETE request to delete current user.
		authenticatedRoutes.DELETE(constants.DeleteCurrentUserPath, func(ginContext *gin.Context) {
			userRouter.userController.DeleteCurrentUser(ginContext)
		})
	}

	// Token-related routes with refresh token authentication middleware.
	tokenRoutes := router.Group("")
	tokenRoutes.Use(httpGinMiddleware.RefreshTokenAuthenticationMiddleware())
	{
		// GET request to refresh access token.
		tokenRoutes.GET(constants.RefreshTokenPath, func(ginContext *gin.Context) {
			userRouter.userController.RefreshAccessToken(ginContext)
		})

		// GET request to logout (invalidate tokens).
		tokenRoutes.GET(constants.LogoutPath, func(ginContext *gin.Context) {
			userRouter.userController.Logout(ginContext)
		})
	}
}
