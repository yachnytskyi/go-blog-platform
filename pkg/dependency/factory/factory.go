package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	zerolog "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger/zerolog"
	useCase "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/usecase"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location            = "pkg.dependency.factory."
	unsupportedLogger   = "Unsupported logger type: %s"
	unsupportedDatabase = "Unsupported database type: %s"
	unsupportedUseCase  = "Unsupported use case type: %s"
	unsupportedDelivery = "Unsupported delivery type: %s"
)

func NewLogger(ctx context.Context) applicationModel.Logger {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Logger {
	case constants.Zerolog:
		return zerolog.NewZerolog()
	// Add other logger options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedLogger, coreConfig.Logger)
		panic(domainError.NewInternalError(location+"NewLoggerFactory", notification))
	}
}

func NewRepositoryFactory(ctx context.Context, logger applicationModel.Logger) applicationModel.Repository {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository(logger)
	// Add other repository options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDatabase, coreConfig.Database)
		logger.Panic(domainError.NewInternalError(location+"NewRepositoryFactory", notification))
		return nil
	}

}

func NewUseCaseFactory(ctx context.Context, logger applicationModel.Logger, repository applicationModel.Repository) applicationModel.UseCase {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.UseCase {
	case constants.UseCase:
		return useCase.NewUseCaseV1(logger)
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedUseCase, coreConfig.UseCase)
		applicationModel.GracefulShutdown(ctx, logger, repository, nil)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", notification))
		return nil
	}
}

func NewDeliveryFactory(ctx context.Context, logger applicationModel.Logger, repository applicationModel.Repository) applicationModel.Delivery {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery(logger)
	// Add other delivery options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDelivery, coreConfig.Delivery)
		applicationModel.GracefulShutdown(ctx, logger, repository, nil)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", notification))
		return nil
	}
}
