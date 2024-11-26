package factory

import (
	"context"
	"fmt"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config"
	configModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	email "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/email"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
)

const (
	location = "pkg.dependency.factory."
)

func NewConfig(configType string) *configModel.ApplicationConfig {
	switch configType {
	case constants.Viper:
		return config.NewViper()
	// Add other config options here as needed.
	default:
		panic(domain.NewInternalError(location+"NewConfig", fmt.Sprintf(constants.UnsupportedConfig, configType)))
	}
}

func NewLogger(config *configModel.ApplicationConfig) interfaces.Logger {
	switch config.Core.Logger {
	case constants.Zerolog:
		return logger.NewZerolog()
	// Add other logger options here as needed.
	default:
		panic(domain.NewInternalError(location+"NewLogger", fmt.Sprintf(constants.UnsupportedLogger, config.Core.Logger)))
	}
}

func NewEmail(config *configModel.ApplicationConfig, logger interfaces.Logger) interfaces.Email {
	switch config.Core.Logger {
	case constants.Zerolog:
		return email.NewGoMail(config, logger)
	// Add other email options here as needed.
	default:
		panic(domain.NewInternalError(location+"NewEmail", fmt.Sprintf(constants.UnsupportedLogger, config.Core.Logger)))
	}
}

func NewRepositoryFactory(config *configModel.ApplicationConfig, logger interfaces.Logger) interfaces.Repository {
	switch config.Core.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository(config, logger)
	// Add other repository options here as needed.
	default:
		logger.Panic(domain.NewInternalError(location+"NewRepositoryFactory", fmt.Sprintf(constants.UnsupportedRepository, config.Core.Database)))
		return nil
	}

}

func NewDeliveryFactory(
	ctx context.Context,
	config *configModel.ApplicationConfig,
	logger interfaces.Logger,
	repository interfaces.Repository) interfaces.Delivery {
	switch config.Core.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery(config, logger)
	// Add other delivery options here as needed.
	default:
		model.GracefulShutdown(ctx, logger, repository)
		logger.Panic(domain.NewInternalError(location+"NewDeliveryFactory", fmt.Sprintf(constants.UnsupportedDelivery, config.Core.Delivery)))
		return nil
	}
}
