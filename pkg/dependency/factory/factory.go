package factory

import (
	"context"
	"fmt"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger"
	useCaseV1 "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/usecase/v1"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location = "pkg.dependency.factory."
)

func NewConfig(configType string) interfaces.Config {
	switch configType {
	case constants.Viper:
		return config.NewViper()
	// Add other logger options here as needed.
	default:
		notification := fmt.Sprintf(model.UnsupportedConfig, configType)
		panic(domainError.NewInternalError(location+"NewLogger", notification))
	}
}

func NewLogger(ctx context.Context, configInstance interfaces.Config) interfaces.Logger {
	config := configInstance.GetConfig()

	switch config.Core.Logger {
	case constants.Zerolog:
		return logger.NewZerolog()
	// Add other logger options here as needed.
	default:
		notification := fmt.Sprintf(model.UnsupportedLogger, config.Core.Logger)
		panic(domainError.NewInternalError(location+"NewLogger", notification))
	}
}

func NewRepositoryFactory(ctx context.Context, configInstance interfaces.Config, logger interfaces.Logger) interfaces.Repository {
	config := configInstance.GetConfig()

	switch config.Core.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository(configInstance, logger)
	// Add other repository options here as needed.
	default:
		notification := fmt.Sprintf(model.UnsupportedRepository, config.Core.Database)
		logger.Panic(domainError.NewInternalError(location+"NewRepositoryFactory", notification))
		return nil
	}

}

func NewUseCaseFactory(ctx context.Context, configInstance interfaces.Config, logger interfaces.Logger, repository interfaces.Repository) interfaces.UseCase {
	config := configInstance.GetConfig()

	switch config.Core.UseCase {
	case constants.UseCaseV1:
		return useCaseV1.NewUseCaseV1(configInstance, logger)
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(model.UnsupportedUseCase, config.Core.UseCase)
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", notification))
		return nil
	}
}

func NewDeliveryFactory(ctx context.Context, configInstance interfaces.Config, logger interfaces.Logger, repository interfaces.Repository) interfaces.Delivery {
	config := configInstance.GetConfig()

	switch config.Core.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery(configInstance, logger)
	// Add other delivery options here as needed.
	default:
		notification := fmt.Sprintf(model.UnsupportedDelivery, config.Core.Delivery)
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", notification))
		return nil
	}
}
