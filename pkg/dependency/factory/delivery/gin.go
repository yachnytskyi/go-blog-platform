package delivery

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/delivery/http/gin"
	postUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin"
	userUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type GinDelivery struct {
	Config *config.ApplicationConfig
	Logger interfaces.Logger
	Server *http.Server
	Router *gin.Engine
}

func NewGinDelivery(config *config.ApplicationConfig, logger interfaces.Logger) *GinDelivery {
	return &GinDelivery{
		Config: config,
		Logger: logger,
	}
}

func (ginDelivery *GinDelivery) CreateDelivery(serverRouters interfaces.ServerRouters) {
	ginDelivery.Router = gin.Default()
	applyMiddleware(ginDelivery.Router, ginDelivery.Config, ginDelivery.Logger)
	configureCORS(ginDelivery.Router, ginDelivery.Config)
	router := ginDelivery.Router.Group(ginDelivery.Config.Gin.ServerGroup)

	// Initialize entity-specific routers.
	serverRouters.UserRouter.Router(router)
	serverRouters.PostRouter.Router(router)

	setNoRouteHandler(ginDelivery.Router, location+"gin.CreateDelivery", ginDelivery.Logger)
	setNoMethodHandler(ginDelivery.Router, location+"gin.CreateDelivery", ginDelivery.Logger)
	ginDelivery.Router.HandleMethodNotAllowed = true

	ginDelivery.Server = &http.Server{
		Addr:    ":" + ginDelivery.Config.Gin.Port,
		Handler: ginDelivery.Router,
	}
}

func (ginDelivery GinDelivery) LaunchServer(ctx context.Context, repository interfaces.Repository) {
	go func() {
		runError := ginDelivery.Router.Run(":" + ginDelivery.Config.Gin.Port)
		if validator.IsError(runError) {
			repository.Close(ctx)
			ginDelivery.Logger.Panic(domain.NewInternalError(location+"gin.LaunchServer.Router.Run", runError.Error()))
		}
	}()

	ginDelivery.Logger.Info(domain.NewInfoMessage(location+"gin.LaunchServer", constants.ServerConnectionSuccess))
}

func (ginDelivery GinDelivery) NewController(useCase any) any {
	switch useCaseType := useCase.(type) {
	case userUseCase.UserUseCase:
		return user.NewUserController(ginDelivery.Config, ginDelivery.Logger, useCaseType)
	case postUseCase.PostUseCase:
		return post.NewPostController(useCaseType)
	default:
		ginDelivery.Logger.Panic(domain.NewInternalError(location+"gin.NewController.default", fmt.Sprintf(constants.UnsupportedDelivery, useCaseType)))
		return nil
	}
}

func (ginDelivery GinDelivery) NewRouter(controller any) interfaces.Router {
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
		ginDelivery.Logger.Panic(domain.NewInternalError(location+"gin.Close.Server.Close", closeError.Error()))
	}

	ginDelivery.Logger.Info(domain.NewInfoMessage(location+"gin.Close", constants.ServerConnectionClosed))
}

func applyMiddleware(router *gin.Engine, config *config.ApplicationConfig, logger interfaces.Logger) {
	router.Use(middleware.RequestIDMiddleware())
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
		httpRequestError := delivery.NewHTTPRequestError(
			location+"gin.setNoRouteHandler.NoRoute",
			requestedPath,
			fmt.Sprintf(constants.RouteNotFoundNotification, requestedPath),
		)
		logger.Error(httpRequestError)
		ginContext.JSON(http.StatusNotFound, httpModel.NewJSONResponseOnFailure(delivery.HandleError(httpRequestError)))
	})
}

func setNoMethodHandler(router *gin.Engine, location string, logger interfaces.Logger) {
	router.NoMethod(func(ginContext *gin.Context) {
		forbiddenMethod := ginContext.Request.Method
		httpRequestError := delivery.NewHTTPRequestError(
			location+"gin.setNoMethodHandler.NoMethod",
			forbiddenMethod,
			fmt.Sprintf(constants.MethodNotAllowedNotification, forbiddenMethod),
		)
		logger.Error(httpRequestError)
		ginContext.JSON(http.StatusMethodNotAllowed, httpModel.NewJSONResponseOnFailure(delivery.HandleError(httpRequestError)))
	})
}
