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

func NewConfig() model.Config {
	switch constants.Config {
	case constants.Viper:
		return config.NewViper()
	// Add other logger options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedLogger, constants.Config)
		panic(domainError.NewInternalError(location+"NewLogger", notification))
	}
}

func NewLogger(ctx context.Context, config model.Config) model.Logger {
	configInstance := config.GetConfig()

	switch configInstance.Core.Logger {
	case constants.Zerolog:
		return logger.NewZerolog()
	// Add other logger options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedLogger, configInstance.Core.Logger)
		panic(domainError.NewInternalError(location+"NewLogger", notification))
	}
}

func NewRepositoryFactory(ctx context.Context, config model.Config, logger model.Logger) model.Repository {
	configInstance := config.GetConfig()

	switch configInstance.Core.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository(config, logger)
	// Add other repository options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDatabase, configInstance.Core.Database)
		logger.Panic(domainError.NewInternalError(location+"NewRepositoryFactory", notification))
		return nil
	}

}

func NewUseCaseFactory(ctx context.Context, config model.Config, logger model.Logger, repository model.Repository) model.UseCase {
	configInstance := config.GetConfig()

	switch configInstance.Core.UseCase {
	case constants.UseCase:
		return useCase.NewUseCaseV1(config, logger)
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedUseCase, configInstance.Core.UseCase)
		model.GracefulShutdown(ctx, logger, repository, nil)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", notification))
		return nil
	}
}

func NewDeliveryFactory(ctx context.Context, config model.Config, logger model.Logger, repository model.Repository) model.Delivery {
	configInstance := config.GetConfig()

	switch configInstance.Core.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery(config, logger)
	// Add other delivery options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDelivery, configInstance.Core.Delivery)
		model.GracefulShutdown(ctx, logger, repository, nil)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", notification))
		return nil
	}
}
