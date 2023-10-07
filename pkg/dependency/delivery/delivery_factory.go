package delivery

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	ginFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/delivery/gin"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location          = "pkg/dependency/delivery/InjectDelivery"
	unsupportedDomain = "unsupported domain type: %s"
)

func InjectDelivery(ctx context.Context, container *applicationModel.Container) {
	applicationConfig := config.AppConfig
	switch applicationConfig.Core.Delivery {
	case constant.Gin:
		container.DeliveryFactory = &ginFactory.GinFactory{Gin: applicationConfig.Gin}
	// Add other domain options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDomain, applicationConfig.Core.Domain)
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Domain:", notification))
		applicationModel.GracefulShutdown(ctx, container)
	}
}
