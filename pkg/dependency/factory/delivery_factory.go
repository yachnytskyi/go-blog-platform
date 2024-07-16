package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	unsupportedDelivery = "Unsupported delivery type: %s"
)

// InjectDelivery injects the appropriate delivery into the container based on the configuration.
// It initializes the delivery according to the application's core configuration settings.
//
// Parameters:
// - ctx: The context for managing request-scoped values, cancellation, and timeouts.
// - container: The dependency injection container where the delivery will be registered.
func InjectDelivery(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration.
	coreConfig := config.GetCoreConfig()

	// Determine the delivery type and inject the corresponding delivery into the container.
	switch coreConfig.Delivery {
	case constants.Gin:
		// Inject the GinDelivery into the container if the delivery type is Gin.
		container.Delivery = delivery.NewGinDelivery()
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
