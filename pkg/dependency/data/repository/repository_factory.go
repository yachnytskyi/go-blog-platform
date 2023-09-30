package repository

import (
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	mongoDBFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository/mongo"
	container "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	application "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/application"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location            = "pkg/dependency/data/repository/InjectRepository"
	unsupportedDatabase = "unsupported database type: %s"
)

func InjectRepository(loadConfig config.Config, container *container.Container) {
	switch loadConfig.Database {
	case config.MongoDB:
		container.RepositoryFactory = &mongoDBFactory.MongoDBFactory{MongoDBConfig: loadConfig.MongoDBConfig}
	// Add other database cases here as needed.
	default:
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Database:", fmt.Sprintf(unsupportedDatabase, loadConfig.Database)))
		application.GracefulShutdown(container)
	}
}