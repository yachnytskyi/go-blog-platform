package factory

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	zerolog "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/logger/zerolog"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
)

const (
	unsupportedLogger = "Unsupported logger type: %s"
)

func NewLoggerFactory(ctx context.Context) applicationModel.Logger {
	coreConfig := config.GetCoreConfig()

	switch coreConfig.Logger {
	case constants.Zerolog:
		return zerolog.NewZerolog()
	// Add other logger options here as needed.
	default:
		notification := fmt.Sprintf(unsupportedLogger, coreConfig.Logger)
		internalError := domainError.NewInternalError(location+"NewLoggerFactory", notification)
		logger.Logger(internalError)
		applicationModel.GracefulShutdown(ctx, nil, nil)
		panic(internalError)
	}
}
