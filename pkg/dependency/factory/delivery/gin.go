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
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	middleware "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/middleware"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "pkg.dependency.delivery.gin."
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

func (ginDelivery GinDelivery) NewUserController(useCase any) interfaces.UserController {
	userUseCase := useCase.(interfaces.UserUseCase)
	return user.NewUserController(ginDelivery.Config, ginDelivery.Logger, userUseCase)
}

func (ginDelivery GinDelivery) NewUserRouter(controller any) interfaces.UserRouter {
	userController := controller.(interfaces.UserController)
	return user.NewUserRouter(ginDelivery.Config, ginDelivery.Logger, userController)
}

func (ginDelivery GinDelivery) NewPostController(userUseCaseInterface, postUseCaseInterface any) interfaces.PostController {
	userUseCase := userUseCaseInterface.(interfaces.UserUseCase)
	postUseCase := postUseCaseInterface.(interfaces.PostUseCase)
	return post.NewPostController(userUseCase, postUseCase)
}

func (ginDelivery GinDelivery) NewPostRouter(controller any) interfaces.PostRouter {
	postController := controller.(interfaces.PostController)
	return post.NewPostRouter(ginDelivery.Config, ginDelivery.Logger, postController)
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
		errorMessage := fmt.Sprintf(constants.RouteNotFoundNotification, requestedPath)
		httpRequestError := httpError.NewHTTPRequestError(location+".setNoRouteHandler.NoRoute", requestedPath, errorMessage)
		logger.Error(httpRequestError)
		ginContext.JSON(constants.StatusNotFound, httpModel.NewJSONResponseOnFailure(httpError.HandleError(httpRequestError)))
	})
}

func setNoMethodHandler(router *gin.Engine, location string, logger interfaces.Logger) {
	router.NoMethod(func(ginContext *gin.Context) {
		forbiddenMethod := ginContext.Request.Method
		errorMessage := fmt.Sprintf(constants.MethodNotAllowedNotification, forbiddenMethod)
		httpRequestError := httpError.NewHTTPRequestError(location+"setNoMethodHandler.NoMethod", forbiddenMethod, errorMessage)
		logger.Error(httpRequestError)
		ginContext.JSON(constants.StatusMethodNotAllowed, httpModel.NewJSONResponseOnFailure(httpError.HandleError(httpRequestError)))
	})
}
