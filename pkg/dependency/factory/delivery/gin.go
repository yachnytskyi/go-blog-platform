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
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "pkg.dependency.delivery.gin."
)

type GinDelivery struct {
	Logger applicationModel.Logger
	Server *http.Server
	Router *gin.Engine
}

func NewGinDelivery(logger applicationModel.Logger) *GinDelivery {
	return &GinDelivery{Logger: logger}
}

func (ginDelivery *GinDelivery) NewDelivery(serverRouters applicationModel.ServerRouters) {
	ginConfig := config.GetGinConfig()
	ginDelivery.Router = gin.Default()
	applyMiddleware(ginDelivery.Router, ginDelivery.Logger)
	configureCORS(ginDelivery.Router, ginConfig)
	router := ginDelivery.Router.Group(ginConfig.ServerGroup)

	// Initialize entity-specific routers.
	serverRouters.UserRouter.UserRouter(router)
	serverRouters.PostRouter.PostRouter(router, serverRouters.UserUseCase)

	setNoRouteHandler(ginDelivery.Router, ginDelivery.Logger)
	setNoMethodHandler(ginDelivery.Router, ginDelivery.Logger)
	ginDelivery.Router.HandleMethodNotAllowed = true

	ginDelivery.Server = &http.Server{
		Addr:    ":" + ginConfig.Port,
		Handler: ginDelivery.Router,
	}
}

func (ginDelivery *GinDelivery) LaunchServer(ctx context.Context, repository applicationModel.Repository) {
	ginConfig := config.GetGinConfig()

	go func() {
		runError := ginDelivery.Router.Run(":" + ginConfig.Port)
		if validator.IsError(runError) {
			repository.CloseRepository(ctx)
			ginDelivery.Logger.Panic(domainError.NewInternalError(location+"LaunchServer.Router.Run", runError.Error()))
		}
	}()

	ginDelivery.Logger.Info(domainError.NewInfoMessage(location+"LaunchServer", constants.ServerConnectionSuccess))
}

func (ginDelivery *GinDelivery) CloseServer(ctx context.Context) {
	shutDownError := ginDelivery.Server.Shutdown(ctx)
	if validator.IsError(shutDownError) {
		ginDelivery.Logger.Panic(domainError.NewInternalError(location+"CloseServer.Server.Shutdown", shutDownError.Error()))
	}

	ginDelivery.Logger.Info(domainError.NewInfoMessage(location+"CloseServer", constants.ServerConnectionClosed))
}

func (ginDelivery *GinDelivery) NewUserController(useCase any) user.UserController {
	userUseCase := useCase.(user.UserUseCase)
	return userDelivery.NewUserController(ginDelivery.Logger, userUseCase)
}

func (ginDelivery *GinDelivery) NewUserRouter(controller any) user.UserRouter {
	userController := controller.(user.UserController)
	return userDelivery.NewUserRouter(ginDelivery.Logger, userController)
}

func (ginDelivery *GinDelivery) NewPostController(userUseCaseInterface, postUseCaseInterface any) post.PostController {
	userUseCase := userUseCaseInterface.(user.UserUseCase)
	postUseCase := postUseCaseInterface.(post.PostUseCase)
	return postDelivery.NewPostController(userUseCase, postUseCase)
}

func (ginDelivery *GinDelivery) NewPostRouter(controller any) post.PostRouter {
	postController := controller.(post.PostController)
	return postDelivery.NewPostRouter(ginDelivery.Logger, postController)
}

func applyMiddleware(router *gin.Engine, logger applicationModel.Logger) {
	router.Use(httpGinMiddleware.CorrelationIDMiddleware())
	router.Use(httpGinMiddleware.SecureHeadersMiddleware())
	router.Use(httpGinMiddleware.CSPMiddleware())
	router.Use(httpGinMiddleware.RateLimitMiddleware())
	router.Use(httpGinMiddleware.ValidateInputMiddleware(logger))
	router.Use(httpGinMiddleware.TimeoutMiddleware(logger))
	router.Use(httpGinMiddleware.LoggerMiddleware(logger))
}

func configureCORS(router *gin.Engine, ginConfig config.Gin) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{ginConfig.AllowOrigins}
	corsConfig.AllowCredentials = ginConfig.AllowCredentials
	router.Use(cors.New(corsConfig))
}

func setNoRouteHandler(router *gin.Engine, logger applicationModel.Logger) {
	router.NoRoute(func(ginContext *gin.Context) {
		requestedPath := ginContext.Request.URL.Path
		errorMessage := fmt.Sprintf(constants.RouteNotFoundNotification, requestedPath)
		httpRequestError := httpError.NewHTTPRequestError(location+"NewDelivery.setNoRouteHandler.ginDelivery.Router.NoRoute", requestedPath, errorMessage)
		logger.Error(httpRequestError)
		httpGinCommon.GinNewJSONFailureResponse(ginContext, httpRequestError, constants.StatusNotFound)
	})
}

func setNoMethodHandler(router *gin.Engine, logger applicationModel.Logger) {
	router.NoMethod(func(ginContext *gin.Context) {
		forbiddenMethod := ginContext.Request.Method
		errorMessage := fmt.Sprintf(constants.MethodNotAllowedNotification, forbiddenMethod)
		httpRequestError := httpError.NewHTTPRequestError(location+"NewDelivery.setNoMethodHandler.ginDelivery.Router.NoMethod", forbiddenMethod, errorMessage)
		logger.Error(httpRequestError)
		httpGinCommon.GinNewJSONFailureResponse(ginContext, httpRequestError, constants.StatusMethodNotAllowed)
	})
}
