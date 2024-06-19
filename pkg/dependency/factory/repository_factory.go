package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	repositoryFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	unsupportedDatabase = "Unsupported database type: %s"
)

// InjectRepository injects the appropriate repository factory into the container based on the configuration.
// It initializes the repository factory according to the application's core configuration settings.
//
// Parameters:
// - ctx: The context for managing request-scoped values, cancellation, and timeouts.
// - container: The dependency injection container where the repository factory will be registered.
func InjectRepository(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration and the MongoDB configuration.
	coreConfig := config.GetCoreConfig()
	mongoDBConfig := config.GetMongoDBConfig()

	// Determine the repository type and inject the corresponding factory into the container.
	switch coreConfig.Database {
	case constants.MongoDB:
		// Inject the MongoDB factory into the container if the database type is MongoDB.
		container.RepositoryFactory = &repositoryFactory.MongoDBFactory{MongoDB: mongoDBConfig}
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
