package model

import (
	"context"

	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	completed = "Graceful shutdown of the app"
)

// GracefulShutdown shuts down the server and closes repository connections gracefully.
//
// This function is responsible for ensuring that all open connections and running servers
// are properly closed before the application exits. It performs the following steps:
// 1. Closes the server using the delivery, if it has been initialized.
// 2. Closes the repository connections using the repository, if it has been initialized.
// 3. Logs the completion of the graceful shutdown process.
//
// Parameters:
// - ctx: The context to control the timeout and cancellation for the shutdown operations.
// - container: The dependency injection container holding the delivery and repository.
func GracefulShutdown(ctx context.Context, container *Container) {
	// Close the server using the delivery, if it is initialized.
	if container.Delivery != nil {
		container.Delivery.CloseServer(ctx)
	}

	// Close the repository connections using the repository, if it is initialized.
	if container.Repository != nil {
		container.Repository.CloseRepository(ctx)
	}

	// Log the completion of the graceful shutdown process.
	logging.Logger(completed)
}
