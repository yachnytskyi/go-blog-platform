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

type GinDelivery struct {
	Server *http.Server // HTTP server instance.
	Router *gin.Engine  // Gin router engine instance.
}

func NewGinDelivery() *GinDelivery {
	return &GinDelivery{}
}

func (ginDelivery *GinDelivery) InitializeServer(serverConfig applicationModel.ServerRouters) {
	ginConfig := config.GetGinConfig()
	ginDelivery.Router = gin.Default()
	applyMiddleware(ginDelivery.Router)
	configureCORS(ginDelivery.Router, ginConfig)
	router := ginDelivery.Router.Group(ginConfig.ServerGroup)

	// Initialize entity-specific routers.
	serverConfig.UserRouter.UserRouter(router)
	serverConfig.PostRouter.PostRouter(router, serverConfig.UserUseCase)

	setNoRouteHandler(ginDelivery.Router)
	setNoMethodHandler(ginDelivery.Router)
	ginDelivery.Router.HandleMethodNotAllowed = true

	ginDelivery.Server = &http.Server{
		Addr:    ":" + ginConfig.Port,
		Handler: ginDelivery.Router,
	}
}

func (ginDelivery *GinDelivery) LaunchServer(ctx context.Context, container *applicationModel.Container) {
	ginConfig := config.GetGinConfig()

	go func() {
		runError := ginDelivery.Router.Run(":" + ginConfig.Port)
		if validator.IsError(runError) {
			container.Repository.CloseRepository(ctx)
			internalError := domainError.NewInternalError(location+"LaunchServer.Router.Run", runError.Error())
			logging.Logger(internalError)
		}
	}()

	logging.Logger(successFullyConnected)
}

func (ginDelivery *GinDelivery) CloseServer(ctx context.Context) {
	shutDownError := ginDelivery.Server.Shutdown(ctx)
	if validator.IsError(shutDownError) {
		internalError := domainError.NewInternalError(location+"CloseServer.Server.Shutdown", shutDownError.Error())
		logging.Logger(internalError)
	}

	logging.Logger(successfullyClosed)
}

func (ginDelivery *GinDelivery) NewUserController(useCase any) user.UserController {
	userUseCase := useCase.(user.UserUseCase)
	return userDelivery.NewUserController(userUseCase)
}

func (ginDelivery *GinDelivery) NewUserRouter(controller any) user.UserRouter {
	userController := controller.(user.UserController)
	return userDelivery.NewUserRouter(userController)
}

func (ginDelivery *GinDelivery) NewPostController(userUseCaseInterface, postUseCaseInterface any) post.PostController {
	userUseCase := userUseCaseInterface.(user.UserUseCase)
	postUseCase := postUseCaseInterface.(post.PostUseCase)
	return postDelivery.NewPostController(userUseCase, postUseCase)
}

func (ginDelivery *GinDelivery) NewPostRouter(controller any) post.PostRouter {
	postController := controller.(post.PostController)
	return postDelivery.NewPostRouter(postController)
}

func applyMiddleware(router *gin.Engine) {
	router.Use(httpGinMiddleware.CorrelationIDMiddleware())
	router.Use(httpGinMiddleware.SecureHeadersMiddleware())
	router.Use(httpGinMiddleware.CSPMiddleware())
	router.Use(httpGinMiddleware.RateLimitMiddleware())
	router.Use(httpGinMiddleware.ValidateInputMiddleware())
	router.Use(httpGinMiddleware.TimeoutMiddleware())
	router.Use(httpGinMiddleware.LoggingMiddleware())
}

func configureCORS(router *gin.Engine, ginConfig config.Gin) {
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{ginConfig.AllowOrigins}
	corsConfig.AllowCredentials = ginConfig.AllowCredentials
	router.Use(cors.New(corsConfig))
}

func setNoRouteHandler(router *gin.Engine) {
	router.NoRoute(func(ginContext *gin.Context) {
		requestedPath := ginContext.Request.URL.Path
		errorMessage := fmt.Sprintf(constants.RouteNotFoundNotification, requestedPath)
		httpRequestError := httpError.NewHTTPRequestError(location+"InitializeServer.setNoRouteHandler.ginDelivery.Router.NoRoute", requestedPath, errorMessage)
		logging.Logger(httpRequestError)
		httpGinCommon.GinNewJSONFailureResponse(ginContext, httpRequestError, constants.StatusNotFound)
	})
}

func setNoMethodHandler(router *gin.Engine) {
	router.NoMethod(func(ginContext *gin.Context) {
		forbiddenMethod := ginContext.Request.Method
		errorMessage := fmt.Sprintf(constants.MethodNotAllowedNotification, forbiddenMethod)
		httpRequestError := httpError.NewHTTPRequestError(location+"InitializeServer.setNoMethodHandler.ginDelivery.Router.NoMethod", forbiddenMethod, errorMessage)
		logging.Logger(httpRequestError)
		httpGinCommon.GinNewJSONFailureResponse(ginContext, httpRequestError, constants.StatusMethodNotAllowed)
	})
}
