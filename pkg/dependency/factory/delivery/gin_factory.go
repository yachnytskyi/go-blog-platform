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
	Gin    config.Gin
	Server *http.Server
	Router *gin.Engine
}

const (
	successfully_connected = "Server is successfully launched..."
	successfully_closed    = "Server has been successfully shutdown..."
)

func (ginFactory *GinFactory) InitializeServer(serverConfig applicationModel.ServerRouters) {
	ginConfig := config.AppConfig.Gin
	ginFactory.Router = gin.Default()
	ginFactory.Router.Use(httpGinMiddleware.TimeoutMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.ValidateInput())
	ginFactory.Router.Use(httpGinMiddleware.SecureHeadersMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.CSPMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.RateLimitMiddleware())
	ginFactory.Router.Use(httpGinMiddleware.LoggingMiddleware())

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{ginConfig.AllowOrigins}
	corsConfig.AllowCredentials = ginConfig.AllowCredentials
	router := ginFactory.Router.Group(ginConfig.ServerGroup)
	ginFactory.Router.Use(cors.New(corsConfig))

	// Routers
	serverConfig.UserRouter.UserRouter(router)
	serverConfig.PostRouter.PostRouter(router)

	ginFactory.Server = &http.Server{
		Addr:    ":" + ginFactory.Gin.Port,
		Handler: ginFactory.Router,
	}
}

func (ginFactory *GinFactory) LaunchServer(ctx context.Context, container *applicationModel.Container) {
	ginConfig := config.AppConfig.Gin

	go func() {
		runError := ginFactory.Router.Run(":" + ginConfig.Port)
		if validator.IsError(runError) {
			container.RepositoryFactory.CloseRepository(ctx)
			runInternalError := domainError.NewInternalError(location+"LaunchServer.Router.Run", runError.Error())
			logging.Logger(runInternalError)
		}
	}()
	logging.Logger(successfully_connected)
}

func (ginFactory *GinFactory) CloseServer(ctx context.Context) {
	shutDownError := ginFactory.Server.Shutdown(ctx)
	if validator.IsError(shutDownError) {
		shutDownInternalError := domainError.NewInternalError(location+"CloseServer.Server.Shutdown", shutDownError.Error())
		logging.Logger(shutDownInternalError)
	}
	logging.Logger(successfully_closed)
}

func (ginFactory *GinFactory) NewUserController(domain any) user.UserController {
	userUseCase := domain.(user.UserUseCase)
	return userDelivery.NewUserController(userUseCase)
}

func (ginFactory *GinFactory) NewUserRouter(controller any) user.UserRouter {
	userController := controller.(user.UserController)
	return userDelivery.NewUserRouter(userController)
}

func (ginFactory *GinFactory) NewPostController(domain any) post.PostController {
	postUseCase := domain.(post.PostUseCase)
	return postDelivery.NewPostController(postUseCase)
}

func (ginFactory *GinFactory) NewPostRouter(controller any) post.PostRouter {
	postController := controller.(post.PostController)
	return postDelivery.NewPostRouter(postController)
}
