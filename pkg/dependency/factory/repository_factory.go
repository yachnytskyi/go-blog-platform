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
	applicationConfig := config.AppConfig
	switch applicationConfig.Core.Database {
	case constants.MongoDB:
		container.RepositoryFactory = &mongoDBFactory.MongoDBFactory{MongoDB: applicationConfig.MongoDB}
	// Add other database cases here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDatabase, applicationConfig.Core.Database)
		logging.Logger(domainError.NewInternalError(location+".applicationConfig.Core.Database:", notification))
		applicationModel.GracefulShutdown(ctx, container)
	}
}
