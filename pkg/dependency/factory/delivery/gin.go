package delivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postDelivery "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userDelivery "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	httpGinMiddleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "pkg.dependency.delivery.gin."
)

type GinDelivery struct {
	Config model.Config
	Logger model.Logger
	Server *http.Server
	Router *gin.Engine
}

func NewGinDelivery(config model.Config, logger model.Logger) *GinDelivery {
	return &GinDelivery{
		Config: config,
		Logger: logger,
	}
}

func (ginDelivery *GinDelivery) NewDelivery(serverRouters model.ServerRouters) {
	config := ginDelivery.Config.GetConfig()
	ginDelivery.Router = gin.Default()
	applyMiddleware(ginDelivery.Router, config, ginDelivery.Logger)
	configureCORS(ginDelivery.Router, config)
	router := ginDelivery.Router.Group(config.Gin.ServerGroup)

	// Initialize entity-specific routers.
	serverRouters.UserRouter.UserRouter(router)
	serverRouters.PostRouter.PostRouter(router, serverRouters.UserUseCase)

	setNoRouteHandler(ginDelivery.Router, location+"NewDelivery", ginDelivery.Logger)
	setNoMethodHandler(ginDelivery.Router, location+"NewDelivery", ginDelivery.Logger)
	ginDelivery.Router.HandleMethodNotAllowed = true

	ginDelivery.Server = &http.Server{
		Addr:    ":" + config.Gin.Port,
		Handler: ginDelivery.Router,
	}
}

func (ginDelivery GinDelivery) LaunchServer(ctx context.Context, repository model.Repository) {
	config := ginDelivery.Config.GetConfig()

	go func() {
		runError := ginDelivery.Router.Run(":" + config.Gin.Port)
		if validator.IsError(runError) {
			repository.Close(ctx)
			ginDelivery.Logger.Panic(domainError.NewInternalError(location+"LaunchServer.Router.Run", runError.Error()))
		}
	}()

	ginDelivery.Logger.Info(domainError.NewInfoMessage(location+"LaunchServer", constants.ServerConnectionSuccess))
}

func (ginDelivery GinDelivery) Close(ctx context.Context) {
	shutDownError := ginDelivery.Server.Shutdown(ctx)
	if validator.IsError(shutDownError) {
		ginDelivery.Logger.Panic(domainError.NewInternalError(location+"Close.Server.Shutdown", shutDownError.Error()))
	}

	ginDelivery.Logger.Info(domainError.NewInfoMessage(location+"Close", constants.ServerConnectionClosed))
}

func (ginDelivery GinDelivery) NewUserController(useCase any) user.UserController {
	userUseCase := useCase.(user.UserUseCase)
	return userDelivery.NewUserController(ginDelivery.Config, ginDelivery.Logger, userUseCase)
}

func (ginDelivery GinDelivery) NewUserRouter(controller any) user.UserRouter {
	userController := controller.(user.UserController)
	return userDelivery.NewUserRouter(ginDelivery.Config, ginDelivery.Logger, userController)
}

func (ginDelivery GinDelivery) NewPostController(userUseCaseInterface, postUseCaseInterface any) post.PostController {
	userUseCase := userUseCaseInterface.(user.UserUseCase)
	postUseCase := postUseCaseInterface.(post.PostUseCase)
	return postDelivery.NewPostController(userUseCase, postUseCase)
}

func (ginDelivery GinDelivery) NewPostRouter(controller any) post.PostRouter {
	postController := controller.(post.PostController)
	return postDelivery.NewPostRouter(ginDelivery.Config, ginDelivery.Logger, postController)
}

func applyMiddleware(router *gin.Engine, config *config.ApplicationConfig, logger model.Logger) {
	router.Use(httpGinMiddleware.CorrelationIDMiddleware())
	router.Use(httpGinMiddleware.SecureHeadersMiddleware(config))
	router.Use(httpGinMiddleware.CSPMiddleware(config))
	router.Use(httpGinMiddleware.RateLimitMiddleware(config))
	router.Use(httpGinMiddleware.ValidateInputMiddleware(config, logger))
	router.Use(httpGinMiddleware.TimeoutMiddleware(logger))
	router.Use(httpGinMiddleware.LoggerMiddleware(logger))
}

func configureCORS(router *gin.Engine, config *config.ApplicationConfig) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Gin.AllowOrigins}
	corsConfig.AllowCredentials = config.Gin.AllowCredentials
	router.Use(cors.New(corsConfig))
}

func setNoRouteHandler(router *gin.Engine, location string, logger model.Logger) {
	router.NoRoute(func(ginContext *gin.Context) {
		requestedPath := ginContext.Request.URL.Path
		errorMessage := fmt.Sprintf(constants.RouteNotFoundNotification, requestedPath)
		httpRequestError := httpError.NewHTTPRequestError(location+".setNoRouteHandler.NoRoute", requestedPath, errorMessage)
		logger.Error(httpRequestError)
		ginContext.JSON(constants.StatusNotFound, httpModel.NewJSONResponseOnFailure(httpError.HandleError(httpRequestError)))
	})
}

func setNoMethodHandler(router *gin.Engine, location string, logger model.Logger) {
	router.NoMethod(func(ginContext *gin.Context) {
		forbiddenMethod := ginContext.Request.Method
		errorMessage := fmt.Sprintf(constants.MethodNotAllowedNotification, forbiddenMethod)
		httpRequestError := httpError.NewHTTPRequestError(location+"setNoMethodHandler.NoMethod", forbiddenMethod, errorMessage)
		logger.Error(httpRequestError)
		ginContext.JSON(constants.StatusMethodNotAllowed, httpModel.NewJSONResponseOnFailure(httpError.HandleError(httpRequestError)))
	})
}
