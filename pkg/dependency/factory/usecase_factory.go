package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	useCase "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/usecase"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location           = "pkg/dependency"
	unsupportedUseCase = "Unsupported use case type: %s"
)

// InjectUseCase injects the appropriate use case into the container based on the configuration.
func InjectUseCase(ctx context.Context, container *applicationModel.Container) {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.UseCase {
	case constants.UseCase:
		container.UseCase = useCase.NewUseCaseV1()
	default:
		notification := fmt.Sprintf(unsupportedUseCase, coreConfig.UseCase)
		internalError := domainError.NewInternalError(location+"InjectUseCase", notification)
		logging.Logger(internalError)
		applicationModel.GracefulShutdown(ctx, container)
		panic(internalError)
	}
}
