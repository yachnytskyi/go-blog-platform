package factory

import (
	"context"
	"fmt"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	email "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/email"
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
		panic(domainError.NewInternalError(location+"NewLogger", fmt.Sprintf(constants.UnsupportedConfig, configType)))
	}
}

func NewLogger(configInstance interfaces.Config) interfaces.Logger {
	config := configInstance.GetConfig()

	switch config.Core.Logger {
	case constants.Zerolog:
		return logger.NewZerolog()
	// Add other logger options here as needed.
	default:
		panic(domainError.NewInternalError(location+"NewLogger", fmt.Sprintf(constants.UnsupportedLogger, config.Core.Logger)))
	}
}

func NewEmail(configInstance interfaces.Config, logger interfaces.Logger) interfaces.Email {
	config := configInstance.GetConfig()

	switch config.Core.Logger {
	case constants.Zerolog:
		return email.NewGoMail(configInstance, logger)
	// Add other logger options here as needed.
	default:
		panic(domainError.NewInternalError(location+"NewLogger", fmt.Sprintf(constants.UnsupportedLogger, config.Core.Logger)))
	}
}

func NewRepositoryFactory(configInstance interfaces.Config, logger interfaces.Logger) interfaces.Repository {
	config := configInstance.GetConfig()

	switch config.Core.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository(configInstance, logger)
	// Add other repository options here as needed.
	default:
		logger.Panic(domainError.NewInternalError(location+"NewRepositoryFactory", fmt.Sprintf(constants.UnsupportedRepository, config.Core.Database)))
		return nil
	}

}

func NewUseCaseFactory(
	ctx context.Context,
	configInstance interfaces.Config,
	logger interfaces.Logger,
	email interfaces.Email,
	repository interfaces.Repository) interfaces.UseCase {
	config := configInstance.GetConfig()

	switch config.Core.UseCase {
	case constants.UseCaseV1:
		return useCaseV1.NewUseCaseV1(configInstance, logger)
	// Add other domain options here as needed.
	default:
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", fmt.Sprintf(constants.UnsupportedUseCase, config.Core.UseCase)))
		return nil
	}
}

func NewDeliveryFactory(
	ctx context.Context,
	configInstance interfaces.Config,
	logger interfaces.Logger,
	repository interfaces.Repository) interfaces.Delivery {
	config := configInstance.GetConfig()

	switch config.Core.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery(configInstance, logger)
	// Add other delivery options here as needed.
	default:
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", fmt.Sprintf(constants.UnsupportedDelivery, config.Core.Delivery)))
		return nil
	}
}
