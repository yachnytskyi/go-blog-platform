package delivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type GinDelivery struct {
	Config interfaces.Config
	Logger interfaces.Logger
	Server *http.Server
	Router *gin.Engine
}

func NewGinDelivery(config interfaces.Config, logger interfaces.Logger) *GinDelivery {
	return &GinDelivery{
		Config: config,
		Logger: logger,
	}
}

func (ginDelivery *GinDelivery) CreateDelivery(serverRouters interfaces.ServerRouters) {
	config := ginDelivery.Config.GetConfig()
	ginDelivery.Router = gin.Default()
	applyMiddleware(ginDelivery.Router, config, ginDelivery.Logger)
	configureCORS(ginDelivery.Router, config)
	router := ginDelivery.Router.Group(config.Gin.ServerGroup)

	// Initialize entity-specific routers.
	serverRouters.UserRouter.UserRouter(router)
	serverRouters.PostRouter.PostRouter(router)

	setNoRouteHandler(ginDelivery.Router, location+"gin.CreateDelivery", ginDelivery.Logger)
	setNoMethodHandler(ginDelivery.Router, location+"gin.CreateDelivery", ginDelivery.Logger)
	ginDelivery.Router.HandleMethodNotAllowed = true

	ginDelivery.Server = &http.Server{
		Addr:    ":" + config.Gin.Port,
		Handler: ginDelivery.Router,
	}
}

func (ginDelivery GinDelivery) LaunchServer(ctx context.Context, repository interfaces.Repository) {
	config := ginDelivery.Config.GetConfig()

	go func() {
		runError := ginDelivery.Router.Run(":" + config.Gin.Port)
		if validator.IsError(runError) {
			repository.Close(ctx)
			ginDelivery.Logger.Panic(domainError.NewInternalError(location+"gin.LaunchServer.Router.Run", runError.Error()))
		}
	}()

	ginDelivery.Logger.Info(domainError.NewInfoMessage(location+"gin.LaunchServer", constants.ServerConnectionSuccess))
}

func (ginDelivery GinDelivery) NewController(userUseCase interfaces.UserUseCase, usecase any) any {
	if usecase == nil {
		return user.NewUserController(ginDelivery.Config, ginDelivery.Logger, userUseCase)
	}

	switch usecaseType := usecase.(type) {
	case interfaces.PostUseCase:
		return post.NewPostController(userUseCase, usecaseType)
	default:
		ginDelivery.Logger.Panic(domainError.NewInternalError(location+"gin.NewController.default", fmt.Sprintf(constants.UnsupportedDelivery, usecaseType)))
		return nil
	}
}

func (ginDelivery GinDelivery) NewRouter(controller any) any {
	switch controllerType := controller.(type) {
	case interfaces.UserController:
		return user.NewUserRouter(ginDelivery.Config, ginDelivery.Logger, controllerType)
	case interfaces.PostController:
		return post.NewPostRouter(ginDelivery.Config, ginDelivery.Logger, controllerType)
	}

	userController := controller.(interfaces.UserController)
	return user.NewUserRouter(ginDelivery.Config, ginDelivery.Logger, userController)
}

func (ginDelivery GinDelivery) Close(ctx context.Context) {
	closeError := ginDelivery.Server.Shutdown(ctx)
	if validator.IsError(closeError) {
		ginDelivery.Logger.Panic(domainError.NewInternalError(location+"gin.Close.Server.Close", closeError.Error()))
	}

	ginDelivery.Logger.Info(domainError.NewInfoMessage(location+"gin.Close", constants.ServerConnectionClosed))
}

func applyMiddleware(router *gin.Engine, config *config.ApplicationConfig, logger interfaces.Logger) {
	router.Use(middleware.CorrelationIDMiddleware())
	router.Use(middleware.SecureHeadersMiddleware(config))
	router.Use(middleware.CSPMiddleware(config))
	router.Use(middleware.RateLimitMiddleware(config))
	router.Use(middleware.ValidateInputMiddleware(config, logger))
	router.Use(middleware.TimeoutMiddleware(logger))
	router.Use(middleware.LoggerMiddleware(logger))
}

func configureCORS(router *gin.Engine, config *config.ApplicationConfig) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{config.Gin.AllowOrigins}
	corsConfig.AllowCredentials = config.Gin.AllowCredentials
	router.Use(cors.New(corsConfig))
}

func setNoRouteHandler(router *gin.Engine, location string, logger interfaces.Logger) {
	router.NoRoute(func(ginContext *gin.Context) {
		requestedPath := ginContext.Request.URL.Path
		httpRequestError := httpError.NewHTTPRequestError(
			location+"gin.setNoRouteHandler.NoRoute",
			requestedPath,
			fmt.Sprintf(constants.RouteNotFoundNotification, requestedPath),
		)
		logger.Error(httpRequestError)
		ginContext.JSON(constants.StatusNotFound, httpModel.NewJSONResponseOnFailure(httpError.HandleError(httpRequestError)))
	})
}

func setNoMethodHandler(router *gin.Engine, location string, logger interfaces.Logger) {
	router.NoMethod(func(ginContext *gin.Context) {
		forbiddenMethod := ginContext.Request.Method
		httpRequestError := httpError.NewHTTPRequestError(
			location+"gin.setNoMethodHandler.NoMethod",
			forbiddenMethod,
			fmt.Sprintf(constants.MethodNotAllowedNotification, forbiddenMethod),
		)
		logger.Error(httpRequestError)
		ginContext.JSON(constants.StatusMethodNotAllowed, httpModel.NewJSONResponseOnFailure(httpError.HandleError(httpRequestError)))
	})
}
