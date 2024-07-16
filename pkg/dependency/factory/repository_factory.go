package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	unsupportedDatabase = "Unsupported database type: %s"
)

// InjectRepository injects the appropriate repository into the container based on the configuration.
// It initializes the concrete repository according to the application's core configuration settings.
//
// Parameters:
// - ctx: The context for managing request-scoped values, cancellation, and timeouts.
// - container: The dependency injection container where the concrete repository will be registered.
func InjectRepository(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration.
	coreConfig := config.GetCoreConfig()

	// Determine the repository type and inject the corresponding repository into the container.
	switch coreConfig.Database {
	case constants.MongoDB:
		// Inject the MongoDB repository into the container if the database type is MongoDB.
		container.Repository = repository.NewMongoDBRepository()
	// Add other database cases here as needed.
	default:
		// Create an error message for the unsupported database type.
		notification := fmt.Sprintf(unsupportedDatabase, coreConfig.Database)
		internalError := domainError.NewInternalError(location+"InjectRepository.applicationConfig.Core.Database:", notification)

		// Log the error.
		logging.Logger(internalError)

		// Perform a graceful shutdown of the application.
		applicationModel.GracefulShutdown(ctx, container)

		// Panic with the error to ensure the application exits.
		panic(internalError)
	}
}
