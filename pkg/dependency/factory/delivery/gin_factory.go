package delivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postDelivery "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userDelivery "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	httpGinCommon "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location              = "pkg.dependency.delivery.gin."
	successFullyConnected = "Server is successfully launched..."
	successfullyClosed    = "Server has been successfully shutdown..."
)

// GinDelivery is responsible for initializing and managing the Gin server instance.
// It holds the HTTP server instance and router engine.
type GinDelivery struct {
	Server *http.Server // HTTP server instance.
	Router *gin.Engine  // Gin router engine instance.
}

// NewGinDelivery creates a new instance of GinDelivery.
//
// Returns:
// - *GinDelivery: The initialized GinDelivery instance.
func NewGinDelivery() *GinDelivery {
	return &GinDelivery{}
}

// InitializeServer sets up the Gin server with the provided routers configuration.
// It loads the Gin configuration, creates a new Gin router engine, applies middleware,
// configures CORS, and initializes entity-specific routers.
//
// Parameters:
// - serverConfig: The configuration for the server routers, including user and post routers.
//
// This method ensures that all necessary middleware and routes are configured
// before the server is launched.
func (ginDelivery *GinDelivery) InitializeServer(serverConfig applicationModel.ServerRouters) {
	// Load the Gin configuration.
	ginConfig := config.GetGinConfig()
	// Create a new Gin router engine instance.
	ginDelivery.Router = gin.Default()

	// Apply middleware to the Gin router.
	applyMiddleware(ginDelivery.Router)

	// Configure CORS settings.
	configureCORS(ginDelivery.Router, ginConfig)

	// Group routes under the server group prefix.
	router := ginDelivery.Router.Group(ginConfig.ServerGroup)

	// Initialize entity-specific routers.
	serverConfig.UserRouter.UserRouter(router)
	serverConfig.PostRouter.PostRouter(router, serverConfig.UserUseCase)

	// Set NoRoute and NoMethod handlers.
	setNoRouteHandler(ginDelivery.Router)
	setNoMethodHandler(ginDelivery.Router)

	// Set HandleMethodNotAllowed.
	ginDelivery.Router.HandleMethodNotAllowed = true

	// Create the HTTP server with the configured Gin router.
	ginDelivery.Server = &http.Server{
		Addr:    ":" + ginConfig.Port,
		Handler: ginDelivery.Router,
	}
}

// LaunchServer starts the Gin server using the provided context and container.
// It runs the Gin router in a separate goroutine and handles any startup errors,
// ensuring proper resource cleanup on failure.
//
// Parameters:
// - ctx: The context for controlling server lifecycle.
// - container: The container holding the application's dependencies.
//
// This method ensures that the server runs asynchronously and logs any errors
// that occur during startup.
func (ginDelivery *GinDelivery) LaunchServer(ctx context.Context, container *applicationModel.Container) {
	ginConfig := config.GetGinConfig()

	go func() {
		// Run the Gin router and handle any errors that occur.
		runError := ginDelivery.Router.Run(":" + ginConfig.Port)
		if validator.IsError(runError) {
			// Close repository on error.
			container.Repository.CloseRepository(ctx)
			// Log the error.
			internalError := domainError.NewInternalError(location+"LaunchServer.Router.Run", runError.Error())
			logging.Logger(internalError)
		}
	}()

	// Log successful server launch.
	logging.Logger(successFullyConnected)
}

// CloseServer gracefully shuts down the server using the provided context.
// It attempts to shutdown the server and logs any errors that occur during the shutdown process.
//
// Parameters:
// - ctx: The context for controlling the shutdown process.
//
// This method ensures that the server is shut down gracefully and logs any issues
// that arise during the shutdown.
func (ginDelivery *GinDelivery) CloseServer(ctx context.Context) {
	// Attempt to shut down the server.
	shutDownError := ginDelivery.Server.Shutdown(ctx)
	if validator.IsError(shutDownError) {
		// Log any errors that occur during shutdown.
		internalError := domainError.NewInternalError(location+"CloseServer.Server.Shutdown", shutDownError.Error())
		logging.Logger(internalError)
	}

	// Log successful server shutdown.
	logging.Logger(successfullyClosed)
}

// NewUserController creates and returns a new UserController instance using the provided domain use case.
//
// Parameters:
// - useCase (any): The use case instance for user operations, must implement user.UserUseCase.
//
// Returns:
// - user.UserController: The newly created UserController instance.
//
// This method type asserts the generic use case to user.UserUseCase and then
// creates a UserController with it.
func (ginDelivery *GinDelivery) NewUserController(useCase any) user.UserController {
	userUseCase := useCase.(user.UserUseCase)
	return userDelivery.NewUserController(userUseCase)
}

// NewUserRouter creates and returns a new UserRouter instance using the provided controller.
//
// Parameters:
// - controller (any): The controller instance for user operations, must implement user.UserController.
//
// Returns:
// - user.UserRouter: The newly created UserRouter instance.
//
// This method type asserts the generic controller to user.UserController and then
// creates a UserRouter with it.
func (ginDelivery *GinDelivery) NewUserRouter(controller any) user.UserRouter {
	userController := controller.(user.UserController)
	return userDelivery.NewUserRouter(userController)
}

// NewPostController creates and returns a new PostController instance using the provided domain use case.
//
// Parameters:
// - userUseCaseInterface (any): The use case instance for user operations, must implement user.UserUseCase.
// - postUseCaseInterface (any): The use case instance for post operations, must implement post.PostUseCase.
//
// Returns:
// - post.PostController: The newly created PostController instance.
//
// This method type asserts the generic use cases to user.UserUseCase and post.PostUseCase
// and then creates a PostController with them.
func (ginDelivery *GinDelivery) NewPostController(userUseCaseInterface, postUseCaseInterface any) post.PostController {
	userUseCase := userUseCaseInterface.(user.UserUseCase)
	postUseCase := postUseCaseInterface.(post.PostUseCase)
	return postDelivery.NewPostController(userUseCase, postUseCase)
}

// NewPostRouter creates and returns a new PostRouter instance using the provided controller.
//
// Parameters:
// - controller (any): The controller instance for post operations, must implement post.PostController.
//
// Returns:
// - post.PostRouter: The newly created PostRouter instance.
//
// This method type asserts the generic controller to post.PostController and then
// creates a PostRouter with it.
func (ginDelivery *GinDelivery) NewPostRouter(controller any) post.PostRouter {
	postController := controller.(post.PostController)
	return postDelivery.NewPostRouter(postController)
}

// Apply middleware to the Gin router.
//
// Parameters:
// - router (*gin.Engine): The Gin router engine to which middleware will be applied.
//
// This function applies a series of middleware functions to the router,
// enhancing its security, logging, and input validation capabilities.
func applyMiddleware(router *gin.Engine) {
	router.Use(httpGinMiddleware.CorrelationIDMiddleware())
	router.Use(httpGinMiddleware.SecureHeadersMiddleware())
	router.Use(httpGinMiddleware.CSPMiddleware())
	router.Use(httpGinMiddleware.RateLimitMiddleware())
	router.Use(httpGinMiddleware.ValidateInputMiddleware())
	router.Use(httpGinMiddleware.TimeoutMiddleware())
	router.Use(httpGinMiddleware.LoggingMiddleware())
}

// Configure CORS settings for the Gin router.
//
// Parameters:
// - router (*gin.Engine): The Gin router engine to configure CORS for.
// - ginConfig (config.Gin): The Gin configuration containing CORS settings.
//
// This function configures the CORS middleware for the router using the
// provided Gin configuration, allowing cross-origin requests from specified origins.
func configureCORS(router *gin.Engine, ginConfig config.Gin) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{ginConfig.AllowOrigins}
	corsConfig.AllowCredentials = ginConfig.AllowCredentials
	router.Use(cors.New(corsConfig))
}

// Set the handler for unmatched routes.
//
// Parameters:
// - router (*gin.Engine): The Gin router engine to set the NoRoute handler for.
//
// This function sets a custom handler for unmatched routes, providing a
// standardized error response for routes that are not found.
func setNoRouteHandler(router *gin.Engine) {
	router.NoRoute(func(ginContext *gin.Context) {
		// Get the requested path that is not found.
		requestedPath := ginContext.Request.URL.Path

		// Create the error message using your constant and the requested path.
		errorMessage := fmt.Sprintf(constants.RouteNotFoundNotification, requestedPath)

		// Create the error with the custom error message.
		httpRequestError := httpError.NewHTTPRequestError(location+"InitializeServer.setNoRouteHandler.ginDelivery.Router.NoRoute", requestedPath, errorMessage)

		// Log the error.
		logging.Logger(httpRequestError)

		// Respond with a not found status and JSON error.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, httpRequestError, constants.StatusNotFound)
	})
}

// Set the handler for methods not allowed.
//
// Parameters:
// - router (*gin.Engine): The Gin router engine to set the NoMethod handler for.
//
// This function sets a custom handler for methods that are not allowed on routes,
// providing a standardized error response for HTTP methods that are not supported.
func setNoMethodHandler(router *gin.Engine) {
	router.NoMethod(func(ginContext *gin.Context) {
		// Get the HTTP method that is not allowed.
		forbiddenMethod := ginContext.Request.Method

		// Create the error message using your constant and the HTTP method.
		errorMessage := fmt.Sprintf(constants.MethodNotAllowedNotification, forbiddenMethod)

		// Create the error with the custom error message.
		httpRequestError := httpError.NewHTTPRequestError(location+"InitializeServer.setNoMethodHandler.ginDelivery.Router.NoMethod", forbiddenMethod, errorMessage)

		// Log the error.
		logging.Logger(httpRequestError)

		// Respond with an unauthorized status and JSON error.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, httpRequestError, constants.StatusMethodNotAllowed)
	})
}
