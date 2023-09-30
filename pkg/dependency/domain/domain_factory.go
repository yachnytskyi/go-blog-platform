package domain

import (
	// "context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	useCaseFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain/usecase"
	container "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/application"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location          = "pkg/dependency/domain/InjectDomain"
	unsupportedDomain = "unsupported domain type: %s"
)

func InjectDomain(loadConfig config.Config, container *container.Container) {
	switch loadConfig.Domain {
	case config.UseCase:
		container.DomainFactory = useCaseFactory.UseCaseFactory{}
	// Add other domain options here as needed.
	default:
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Domain:", fmt.Sprintf(unsupportedDomain, loadConfig.Domain)))
		application.GracefulShutdown(container)
	}
}
