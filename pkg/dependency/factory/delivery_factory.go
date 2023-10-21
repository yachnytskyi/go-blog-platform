package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	ginFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery/gin"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	unsupportedDelivery = "Unsupported delivery type: %s"
)

func InjectDelivery(ctx context.Context, container *applicationModel.Container) {
	applicationConfig := config.AppConfig
	switch applicationConfig.Core.Delivery {
	case constants.Gin:
		container.DeliveryFactory = &ginFactory.GinFactory{Gin: applicationConfig.Gin}
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDelivery, applicationConfig.Core.Delivery)
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Delivery:", notification))
		applicationModel.GracefulShutdown(ctx, container)
	}
}
