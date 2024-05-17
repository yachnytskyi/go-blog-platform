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

// InjectDomain injects the appropriate domain factory into the container based on the configuration.
func InjectDomain(ctx context.Context, container *applicationModel.Container) {
	// Load the core configuration.
	coreConfig := config.AppConfig.Core

	// Switch based on the configured domain type.
	switch coreConfig.Domain {
	case constants.UseCase:
		// Inject the UseCaseFactory into the container if the domain type is UseCase.
		container.DomainFactory = domainFactory.UseCaseFactory{}
	// Add other domain options here as needed.
	default:
		// Handle unsupported domain types by logging an error and shutting down the application gracefully.
		notification := fmt.Sprintf(unsupportedDomain, coreConfig.Domain)
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Domain:", notification))
		applicationModel.GracefulShutdown(ctx, container)
	}
}
