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
	useCase "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/usecase"
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
		panic(domainError.NewInternalError(location+"NewConfig", fmt.Sprintf(constants.UnsupportedConfig, configType)))
	}
}

func NewLogger(config interfaces.Config) interfaces.Logger {
	configInstance := config.GetConfig()

	switch configInstance.Core.Logger {
	case constants.Zerolog:
		return logger.NewZerolog()
	// Add other logger options here as needed.
	default:
		panic(domainError.NewInternalError(location+"NewLogger", fmt.Sprintf(constants.UnsupportedLogger, configInstance.Core.Logger)))
	}
}

func NewEmail(config interfaces.Config, logger interfaces.Logger) interfaces.Email {
	configInstance := config.GetConfig()

	switch configInstance.Core.Logger {
	case constants.Zerolog:
		return email.NewGoMail(config, logger)
	// Add other logger options here as needed.
	default:
		panic(domainError.NewInternalError(location+"NewEmail", fmt.Sprintf(constants.UnsupportedLogger, configInstance.Core.Logger)))
	}
}

func NewRepositoryFactory(config interfaces.Config, logger interfaces.Logger) interfaces.Repository {
	configInstance := config.GetConfig()

	switch configInstance.Core.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository(config, logger)
	// Add other repository options here as needed.
	default:
		logger.Panic(domainError.NewInternalError(location+"NewRepositoryFactory", fmt.Sprintf(constants.UnsupportedRepository, configInstance.Core.Database)))
		return nil
	}

}

func NewUseCaseFactory(
	ctx context.Context,
	config interfaces.Config,
	logger interfaces.Logger,
	email interfaces.Email,
	repository interfaces.Repository) interfaces.UseCase {

	configInstance := config.GetConfig()
	switch configInstance.Core.UseCase {
	case constants.UseCaseV1:
		return useCase.NewUseCaseV1(config, logger)
	// Add other domain options here as needed.
	default:
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domainError.NewInternalError(location+"NewUseCaseFactory", fmt.Sprintf(constants.UnsupportedUseCase, configInstance.Core.UseCase)))
		return nil
	}
}

func NewDeliveryFactory(
	ctx context.Context,
	config interfaces.Config,
	logger interfaces.Logger,
	repository interfaces.Repository) interfaces.Delivery {

	configInstance := config.GetConfig()
	switch configInstance.Core.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery(config, logger)
	// Add other delivery options here as needed.
	default:
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domainError.NewInternalError(location+"NewDeliveryFactory", fmt.Sprintf(constants.UnsupportedDelivery, configInstance.Core.Delivery)))
		return nil
	}
}
