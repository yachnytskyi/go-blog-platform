package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	useCaseFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/usecase"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	unsupportedUseCase = "unsupported usecase type: %s"
)

// InjectUseCase injects the appropriate usecase factory into the container based on the configuration.
// It initializes the use case factory according to the application's core configuration settings.
//
// Parameters:
// - ctx: The context for managing request-scoped values, cancellation, and timeouts.
// - container: The dependency injection container where the use case factory will be registered.
func InjectUseCase(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration from the application's configuration.
	coreConfig := config.AppConfig.Core

	// Switch based on the configured use case type.
	switch coreConfig.UseCase {
	case constants.UseCase:
		// Inject the UseCaseFactory into the container if the usecase type is UseCaseV1.
		container.UseCaseFactory = useCaseFactory.UseCaseFactoryV1{}
	// Add other use case options here as needed.
	default:
		// Handle unsupported usecase types by logging an error and shutting down the application gracefully.
		notification := fmt.Sprintf(unsupportedUseCase, coreConfig.UseCase)
		logging.Logger(domainError.NewInternalError(location+".loadConfig.UseCase:", notification))
		applicationModel.GracefulShutdown(ctx, container)
		panic(notification)
	}
}
