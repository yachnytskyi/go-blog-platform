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

func InjectRepository(ctx context.Context, container *applicationModel.Container) {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Database {
	case constants.MongoDB:
		container.Repository = repository.NewMongoDBRepository()
	default:
		notification := fmt.Sprintf(unsupportedDatabase, coreConfig.Database)
		internalError := domainError.NewInternalError(location+"InjectRepository", notification)
		logging.Logger(internalError)
		applicationModel.GracefulShutdown(ctx, container)
		panic(internalError)
	}
}
