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

// InjectDelivery injects the appropriate delivery factory into the container based on the configuration.
func InjectDelivery(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration and the Gin configuration.
	coreConfig := config.AppConfig.Core
	ginConfig := config.AppConfig.Gin

	// Switch based on the configured delivery type.
	switch coreConfig.Delivery {
	case constants.Gin:
		// Inject the GinFactory into the container if the delivery type is Gin.
		container.DeliveryFactory = &deliveryFactory.GinFactory{Gin: ginConfig}
	// Add other delivery options here as needed.
	default:
		// Handle unsupported delivery types by logging an error and shutting down the application gracefully.
		notification := fmt.Sprintf(unsupportedDelivery, coreConfig.Delivery)
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Delivery:", notification))
		applicationModel.GracefulShutdown(ctx, container)
		panic(notification)
	}
}
