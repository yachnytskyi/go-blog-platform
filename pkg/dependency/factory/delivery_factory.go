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
// It initializes the delivery factory according to the application's core configuration settings.
//
// Parameters:
// - ctx: The context for managing request-scoped values, cancellation, and timeouts.
// - container: The dependency injection container where the delivery factory will be registered.
func InjectDelivery(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration and the Gin configuration.
	coreConfig := config.GetCoreConfig()
	ginConfig := config.GetGinConfig()

	// Determine the delivery type and inject the corresponding factory into the container.
	switch coreConfig.Delivery {
	case constants.Gin:
		// Inject the GinFactory into the container if the delivery type is Gin.
		container.DeliveryFactory = &deliveryFactory.GinFactory{Gin: ginConfig}
		// Add other delivery options here as needed.
	default:
		// Create an error message for the unsupported delivery type.
		notification := fmt.Sprintf(unsupportedDelivery, coreConfig.Delivery)
		internalError := domainError.NewInternalError(location+".InjectDelivery.loadConfig.Delivery:", notification)

		// Log the error.
		logging.Logger(internalError)

		// Perform a graceful shutdown of the application.
		applicationModel.GracefulShutdown(ctx, container)

		// Panic with the error to ensure the application exits.
		panic(internalError)
	}
}
