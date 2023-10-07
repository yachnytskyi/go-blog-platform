package model

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	completed = "Completed graceful shutdown of the app"
)

func GracefulShutdown(ctx context.Context, container *Container) {
	container.DeliveryFactory.CloseServer(ctx)
	container.RepositoryFactory.CloseRepository(ctx)
	logging.Logger(completed)
}
