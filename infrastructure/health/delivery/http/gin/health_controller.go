package gin

import (
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
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
	ginContext.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
		"details": gin.H{
			"database": dbHealthy,
		},
	})
}
