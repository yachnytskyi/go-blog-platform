package model

import (
	"context"

	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	completed = "Graceful shutdown of the app"
)

// GracefulShutdown shuts down the server and closes repository connections gracefully.
func GracefulShutdown(ctx context.Context, container *Container) {
	// Close the server using the delivery factory.
	container.DeliveryFactory.CloseServer(ctx)
	// Close the repository connections using the repository factory.
	container.RepositoryFactory.CloseRepository(ctx)
	// Log the completion of the graceful shutdown process.
	logging.Logger(completed)
}
