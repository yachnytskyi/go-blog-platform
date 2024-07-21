package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	useCase "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/usecase"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
)

const (
	location           = "pkg.dependency.factory."
	unsupportedUseCase = "Unsupported use case type: %s"
)

func NewUseCaseFactory(ctx context.Context, repository applicationModel.Repository) applicationModel.UseCase {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.UseCase {
	case constants.UseCase:
		return useCase.NewUseCaseV1()
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedUseCase, coreConfig.UseCase)
		internalError := domainError.NewInternalError(location+"NewUseCaseFactory", notification)
		logger.Logger(internalError)
		applicationModel.GracefulShutdown(ctx, repository, nil)
		panic(internalError)
	}
}
