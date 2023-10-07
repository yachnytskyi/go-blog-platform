package model

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	completed = "Completed graceful shutdown of the app"
)

func GracefulShutdown(ctx context.Context, container *Container) {
	container.DeliveryFactory.CloseServer()
	container.RepositoryFactory.CloseRepository()
	logging.Logger(completed)
}
