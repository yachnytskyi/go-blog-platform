package model

import (
	"context"

	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	completed = "Graceful shutdown of the app"
)

func GracefulShutdown(ctx context.Context, container *Container) {
	container.DeliveryFactory.CloseServer(ctx)
	container.RepositoryFactory.CloseRepository(ctx)
	logging.Logger(completed)
}
