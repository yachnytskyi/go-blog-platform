package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	useCase "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/usecase"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location           = "pkg/dependency"
	unsupportedUseCase = "Unsupported use case type: %s"
)

// InjectUseCase injects the appropriate use case into the container based on the configuration.
// It initializes the use case according to the application's core configuration settings.
//
// Parameters:
// - ctx: The context for managing request-scoped values, cancellation, and timeouts.
// - container: The dependency injection container where the use case will be registered.
func InjectUseCase(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration from the application's configuration.
	coreConfig := config.GetCoreConfig()

	// Determine the use case type and inject the corresponding use case into the container.
	switch coreConfig.UseCase {
	case constants.UseCase:
		// Inject the UseCaseV1 into the container if the use case type is UseCaseV1.
		container.UseCase = useCase.NewUseCaseV1()
	// Add other use case options here as needed.
	default:
		// Create an error message for the unsupported use case type.
		notification := fmt.Sprintf(unsupportedUseCase, coreConfig.UseCase)
		internalError := domainError.NewInternalError(location+"InjectUseCase.loadConfig.UseCase:", notification)

		// Log the error.
		logging.Logger(internalError)

		// Perform a graceful shutdown of the application.
		applicationModel.GracefulShutdown(ctx, container)

		// Panic with the error to ensure the application exits.
		panic(internalError)
	}
}
