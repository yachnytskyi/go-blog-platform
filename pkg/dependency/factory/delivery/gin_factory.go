package delivery

import (
	"context"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postDelivery "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userDelivery "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "pkg.dependency.delivery.gin."
)

type GinFactory struct {
	Gin    config.Gin   // Gin configuration.
	Server *http.Server // HTTP server instance.
	Router *gin.Engine  // Gin router engine instance.
}

const (
	successFullyConnected = "Server is successfully launched..."
	successfullyClosed    = "Server has been successfully shutdown..."
)

// InitializeServer sets up the Gin server with the provided routers configuration.
// It loads the Gin configuration, creates a new Gin router engine, applies middleware,
// configures CORS, and initializes entity-specific routers.
func (ginFactory *GinFactory) InitializeServer(serverConfig applicationModel.ServerRouters) {
	// Load the Gin configuration.
	ginConfig := config.AppConfig.Gin
	// Create a new Gin router engine instance.
	ginFactory.Router = gin.Default()
	// Apply middleware to the Gin router.
	ginFactory.Router.Use(httpGinMiddleware.TimeoutMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.ValidateInputMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.SecureHeadersMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.CSPMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.RateLimitMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.LoggingMiddleware())

	// Configure CORS settings.
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{ginConfig.AllowOrigins}
	corsConfig.AllowCredentials = ginConfig.AllowCredentials
	ginFactory.Router.Use(cors.New(corsConfig))

	// Group routes under the server group prefix.
	router := ginFactory.Router.Group(ginConfig.ServerGroup)

	// Initialize entity-specific routers.
	serverConfig.UserRouter.UserRouter(router)
	serverConfig.PostRouter.PostRouter(router)

	// Create the HTTP server with the configured Gin router.
	ginFactory.Server = &http.Server{
		Addr:    ":" + ginFactory.Gin.Port,
		Handler: ginFactory.Router,
	}
}

// LaunchServer starts the Gin server using the provided context and container.
// It runs the Gin router in a separate goroutine and handles any startup errors,
// ensuring proper resource cleanup on failure.
func (ginFactory *GinFactory) LaunchServer(ctx context.Context, container *applicationModel.Container) {
	ginConfig := config.AppConfig.Gin

	go func() {
		// Run the Gin router and handle any errors that occur.
		runError := ginFactory.Router.Run(":" + ginConfig.Port)
		if validator.IsError(runError) {
			// Close repository on error.
			container.RepositoryFactory.CloseRepository(ctx)
			// Log the error.
			runInternalError := domainError.NewInternalError(location+"LaunchServer.Router.Run", runError.Error())
			logging.Logger(runInternalError)
		}
	}()

	// Log successful server launch.
	logging.Logger(successFullyConnected)
}

// CloseServer gracefully shuts down the server using the provided context.
// It attempts to shutdown the server and logs any errors that occur during the shutdown process.
func (ginFactory *GinFactory) CloseServer(ctx context.Context) {
	// Attempt to shut down the server.
	shutDownError := ginFactory.Server.Shutdown(ctx)
	if validator.IsError(shutDownError) {
		// Log any errors that occur during shutdown.
		shutDownInternalError := domainError.NewInternalError(location+"CloseServer.Server.Shutdown", shutDownError.Error())
		logging.Logger(shutDownInternalError)
	}
	// Log successful server shutdown.
	logging.Logger(successfullyClosed)
}

// NewUserController creates and returns a new UserController instance using the provided domain use case.
// It casts the generic use case parameter to the specific user.UserUseCase type and creates the UserController.
func (ginFactory *GinFactory) NewUserController(useCase any) user.UserController {
	userUseCase := useCase.(user.UserUseCase)
	return userDelivery.NewUserController(userUseCase)
}

// NewUserRouter creates and returns a new UserRouter instance using the provided controller.
// It casts the generic controller parameter to the specific user.UserController type and creates the UserRouter.
func (ginFactory *GinFactory) NewUserRouter(controller any) user.UserRouter {
	userController := controller.(user.UserController)
	return userDelivery.NewUserRouter(userController)
}

// NewPostController creates and returns a new PostController instance using the provided domain use case.
// It casts the generic use case parameter to the specific post.PostUseCase type and creates the PostController.
func (ginFactory *GinFactory) NewPostController(useCase any) post.PostController {
	postUseCase := useCase.(post.PostUseCase)
	return postDelivery.NewPostController(postUseCase)
}

// NewPostRouter creates and returns a new PostRouter instance using the provided controller.
// It casts the generic controller parameter to the specific post.PostController type and creates the PostRouter.
func (ginFactory *GinFactory) NewPostRouter(controller any) post.PostRouter {
	postController := controller.(post.PostController)
	return postDelivery.NewPostRouter(postController)
}
