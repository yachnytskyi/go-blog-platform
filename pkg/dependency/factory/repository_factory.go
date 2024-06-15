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
	location            = "pkg/dependency/data/repository/InjectRepository."
	unsupportedDatabase = "Unsupported database type: %s"
)

// InjectRepository injects the appropriate repository factory into the container based on the configuration.
func InjectRepository(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration and the MongoDB configuration.
	coreConfig := config.AppConfig.Core
	mongoDBConfig := config.AppConfig.MongoDB

	// Switch based on the configured database type.
	switch coreConfig.Database {
	case constants.MongoDB:
		// Inject the MongoDB factory into the container if the database type is MongoDB.
		container.RepositoryFactory = &repositoryFactory.MongoDBFactory{MongoDB: mongoDBConfig}
	// Add other database cases here as needed.
	default:
		// Handle unsupported database types by logging an error and shutting down the application gracefully.
		notification := fmt.Sprintf(unsupportedDatabase, coreConfig.Database)
		internalError := domainError.NewInternalError(location+".applicationConfig.Core.Database:", notification)
		logging.Logger(internalError)
		applicationModel.GracefulShutdown(ctx, container)
		panic(internalError)
	}
}
