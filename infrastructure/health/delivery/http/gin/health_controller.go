package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	view "github.com/yachnytskyi/golang-mongo-grpc/infrastructure/health/delivery/http/model"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
)

type HealthCheckController struct {
	Config     *config.ApplicationConfig
	Logger     interfaces.Logger
	Repository interfaces.Repository
}

func NewHealthCheckController(config *config.ApplicationConfig, logger interfaces.Logger, repository interfaces.Repository) HealthCheckController {
	return HealthCheckController{
		Config:     config,
		Logger:     logger,
		Repository: repository,
	}
}

func (healthCheckController HealthCheckController) HealthCheck(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)

	dbHealthy := healthCheckController.Repository.DatabasePing()
	if !dbHealthy {
		ginContext.JSON(http.StatusServiceUnavailable, model.JSONResponseOnSuccess{Data: view.NewHealthStatus(dbHealthy), Status: constants.Fail})
		return
	}

	ginContext.JSON(http.StatusOK, model.NewJSONResponseOnSuccess(view.NewHealthStatus(dbHealthy)))
}
