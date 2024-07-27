package factory

import (
	"context"
	"fmt"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger"
	useCase "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/usecase"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location            = "pkg.dependency.factory."
	unsupportedConfig   = "Unsupported Config type: %s"
	unsupportedLogger   = "Unsupported logger type: %s"
	unsupportedDatabase = "Unsupported database type: %s"
	unsupportedUseCase  = "Unsupported use case type: %s"
	unsupportedDelivery = "Unsupported delivery type: %s"
)

func NewConfig(configType string) model.Config {
	switch configType {
	case constants.Viper:
		return config.NewViper()
	// Add other logger options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedConfig, configType)
		panic(domainError.NewInternalError(location+"NewLogger", notification))
	}
}

func NewLogger(ctx context.Context, configInstance model.Config) model.Logger {
	config := configInstance.GetConfig()

	switch config.Core.Logger {
	case constants.Zerolog:
		return logger.NewZerolog()
	// Add other logger options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedLogger, config.Core.Logger)
		panic(domainError.NewInternalError(location+"NewLogger", notification))
	}
}

func NewRepositoryFactory(ctx context.Context, configInstance model.Config, logger model.Logger) model.Repository {
	config := configInstance.GetConfig()

	switch config.Core.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository(configInstance, logger)
	// Add other repository options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDatabase, config.Core.Database)
		logger.Panic(domainError.NewInternalError(location+"NewRepositoryFactory", notification))
		return nil
	}

}

func NewUseCaseFactory(ctx context.Context, configInstance model.Config, logger model.Logger, repository model.Repository) model.UseCase {
	config := configInstance.GetConfig()

	switch config.Core.UseCase {
	case constants.UseCase:
		return useCase.NewUseCaseV1(configInstance, logger)
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedUseCase, config.Core.UseCase)
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", notification))
		return nil
	}
}

func NewDeliveryFactory(ctx context.Context, configInstance model.Config, logger model.Logger, repository model.Repository) model.Delivery {
	config := configInstance.GetConfig()

	switch config.Core.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery(configInstance, logger)
	// Add other delivery options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDelivery, config.Core.Delivery)
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", notification))
		return nil
	}
}
