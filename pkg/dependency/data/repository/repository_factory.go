package repository

import (
	"fmt"

	"github.com/yachnytskyi/golang-mongo-grpc/config"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	mongoDBFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository/mongo"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	application "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/application"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location            = "pkg/dependency/data/repository/InjectRepository"
	unsupportedDatabase = "unsupported database type: %s"
)

func InjectRepository(container *applicationModel.Container) {
	applicationConfig := config.AppConfig
	switch applicationConfig.Database {
	case constant.MongoDB:
		container.RepositoryFactory = &mongoDBFactory.MongoDBFactory{MongoDBConfig: applicationConfig.MongoDBConfig}
	// Add other database cases here as needed.
	default:
		logging.Logger(domainError.NewInternalError(location+".applicationConfig.Database:", fmt.Sprintf(unsupportedDatabase, applicationConfig.Database)))
		application.GracefulShutdown(container)
	}
}
