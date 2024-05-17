package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/domain"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	unsupportedDomain = "unsupported domain type: %s"
)

func InjectDomain(ctx context.Context, container *applicationModel.Container) {
	coreConfig := config.AppConfig.Core

	switch coreConfig.Domain {
	case constants.UseCase:
		container.DomainFactory = domainFactory.UseCaseFactory{}
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDomain, coreConfig.Domain)
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Domain:", notification))
		applicationModel.GracefulShutdown(ctx, container)
	}
}
