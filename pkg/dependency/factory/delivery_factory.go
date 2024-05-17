package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	deliveryFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	unsupportedDelivery = "Unsupported delivery type: %s"
)

func InjectDelivery(ctx context.Context, container *applicationModel.Container) {
	coreConfig := config.AppConfig.Core
	ginConfig := config.AppConfig.Gin

	switch coreConfig.Delivery {
	case constants.Gin:
		container.DeliveryFactory = &deliveryFactory.GinFactory{Gin: ginConfig}
	// Add other delivery options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDelivery, coreConfig.Delivery)
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Delivery:", notification))
		applicationModel.GracefulShutdown(ctx, container)
	}
}
