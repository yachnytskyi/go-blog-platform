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

func NewDeliveryFactory(ctx context.Context, repository applicationModel.Repository) applicationModel.Delivery {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Delivery {
	case constants.Gin:
		return delivery.NewGinDelivery()
	// Add other delivery options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedDelivery, coreConfig.Delivery)
		internalError := domainError.NewInternalError(location+".NewDeliveryFactory", notification)
		logger.Logger(internalError)
		applicationModel.GracefulShutdown(ctx, repository, nil)
		panic(internalError)
	}
}
