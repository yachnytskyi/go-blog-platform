package model

import (
	"context"

	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
)

const (
	completed = "Graceful shutdown of the app"
)

func GracefulShutdown(ctx context.Context, container *Container) {
	if container.Delivery != nil {
		container.Delivery.CloseServer(ctx)
	}
	if container.Repository != nil {
		container.Repository.CloseRepository(ctx)
	}

	logger.Logger(completed)
}
