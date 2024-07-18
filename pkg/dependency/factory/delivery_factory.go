package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/delivery"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
)

const (
	unsupportedDelivery = "Unsupported delivery type: %s"
)

func InjectDelivery(ctx context.Context, container *applicationModel.Container) {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Delivery {
	case constants.Gin:
		container.Delivery = delivery.NewGinDelivery()
	default:
		notification := fmt.Sprintf(unsupportedDelivery, coreConfig.Delivery)
		internalError := domainError.NewInternalError(location+".InjectDelivery", notification)
		logger.Logger(internalError)
		applicationModel.GracefulShutdown(ctx, container)
		panic(internalError)
	}
}
