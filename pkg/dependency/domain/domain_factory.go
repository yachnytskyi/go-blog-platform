package domain

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	useCaseFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain/usecase"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location          = "pkg/dependency/domain/InjectDomain."
	unsupportedDomain = "unsupported domain type: %s"
)

func InjectDomain(ctx context.Context, container *applicationModel.Container) {
	applicationConfig := config.AppConfig
	switch applicationConfig.Core.Domain {
	case constants.UseCase:
		container.DomainFactory = useCaseFactory.UseCaseFactory{}
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDomain, applicationConfig.Core.Domain)
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Domain:", notification))
		applicationModel.GracefulShutdown(ctx, container)
	}
}
