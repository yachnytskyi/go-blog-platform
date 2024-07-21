package model

import (
	"context"

	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
)

const (
	completed = "Graceful shutdown of the app"
)

func GracefulShutdown(ctx context.Context, repository Repository, delivery Delivery) {
	if delivery != nil {
		delivery.CloseServer(ctx)
	}
	if repository != nil {
		repository.CloseRepository(ctx)
	}

	logger.Logger(completed)
}
