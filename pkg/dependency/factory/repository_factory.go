package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	mongoDBFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository/mongo"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location            = "pkg/dependency/data/repository/InjectRepository."
	unsupportedDatabase = "Unsupported database type: %s"
)

func InjectRepository(ctx context.Context, container *applicationModel.Container) {
	coreConfig := config.AppConfig.Core
	mongoDBConfig := config.AppConfig.MongoDB

	switch coreConfig.Database {
	case constants.MongoDB:
		container.RepositoryFactory = &mongoDBFactory.MongoDBFactory{MongoDB: mongoDBConfig}
	// Add other database cases here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDatabase, coreConfig.Database)
		logging.Logger(domainError.NewInternalError(location+".applicationConfig.Core.Database:", notification))
		applicationModel.GracefulShutdown(ctx, container)
	}
}
