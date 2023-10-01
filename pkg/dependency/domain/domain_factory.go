package domain

import (
	// "context"
	"fmt"

	"github.com/yachnytskyi/golang-mongo-grpc/config"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
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

func InjectDomain(loadConfig config.ApplicationConfig, container *container.Container) {
	switch loadConfig.Domain {
	case constant.UseCase:
		container.DomainFactory = useCaseFactory.UseCaseFactory{}
	// Add other domain options here as needed.
	default:
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Domain:", fmt.Sprintf(unsupportedDomain, loadConfig.Domain)))
		application.GracefulShutdown(container)
	}
}
