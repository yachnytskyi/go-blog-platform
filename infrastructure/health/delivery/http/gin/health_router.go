package gin

import (
	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
)

type HealthCheckRouter struct {
	Config                *config.ApplicationConfig
	Logger                interfaces.Logger
	HealthCheckController interfaces.HealthCheckController
	Repository            interfaces.Repository
}

func NewHealthCheckRouter(config *config.ApplicationConfig, logger interfaces.Logger, heathCheckController interfaces.HealthCheckController, repository interfaces.Repository) HealthCheckRouter {
	return HealthCheckRouter{
		Config:                config,
		Logger:                logger,
		HealthCheckController: heathCheckController,
		Repository:            repository,
	}
}

// HealthCheckRouter defines the health-related routes and connects them to the corresponding controller methods.
func (healthCheckRouter HealthCheckRouter) Router(routerGroup any) {
	ginRouterGroup := routerGroup.(*gin.RouterGroup)
	router := ginRouterGroup.Group("/health")

	// Public routes.
	publicRoutes := router.Group("")
	{
		publicRoutes.GET("", func(ginContext *gin.Context) {
			healthCheckRouter.HealthCheckController.HealthCheck(ginContext)
		})
	}
}
