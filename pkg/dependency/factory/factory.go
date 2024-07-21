package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/data/repository"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger/zerolog"
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

func NewLoggerFactory(ctx context.Context) applicationModel.Logger {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Logger {
	case constants.Zerolog:
		return zerolog.NewZerolog()
	// Add other logger options here as needed.
	default:
		applicationModel.GracefulShutdown(ctx, nil, nil)
		notification := fmt.Sprintf(unsupportedLogger, coreConfig.Logger)
		panic(domainError.NewInternalError(location+"NewLoggerFactory", notification))
	}
}

func NewRepositoryFactory(ctx context.Context) applicationModel.Repository {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Database {
	case constants.MongoDB:
		return repository.NewMongoDBRepository()
	// Add other repository options here as needed.
	default:
		panic(handleError(ctx, location+"NewRepositoryFactory", unsupportedDatabase, coreConfig.UseCase, nil))
	}
}

func NewUseCaseFactory(ctx context.Context, repository applicationModel.Repository) applicationModel.UseCase {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.UseCase {
	case constants.UseCase:
		return useCase.NewUseCaseV1()
	// Add other domain options here as needed.
	default:
		panic(handleError(ctx, location+"NewUseCaseFactory", unsupportedUseCase, coreConfig.UseCase, repository))
	}
}

func NewDeliveryFactory(ctx context.Context, repository applicationModel.Repository) applicationModel.Delivery {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery()
	// Add other delivery options here as needed.
	default:
		panic(handleError(ctx, location+"NewUseCaseFactory", unsupportedUseCase, coreConfig.UseCase, repository))
	}
}

// handleError handles unsupported configuration options, logs and returns the error.
func handleError(ctx context.Context, location, format, value string, repository applicationModel.Repository) error {
	notification := fmt.Sprintf(format, value)
	internalError := domainError.NewInternalError(location+".handleError", notification)
	applicationModel.GracefulShutdown(ctx, repository, nil)
	return internalError
}
