package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	dependency "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	// Message to log if the application fails to initialize.
	failedInitialization = "Failed to initialize application"
)

// init is called before the main function to perform setup tasks.
func init() {
	// Load configuration settings from config files or environment variables.
	config.LoadConfig()
}

func main() {
	// Create a context with a timeout for the entire application lifecycle.
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimer)
	defer cancel() // Ensure the context is canceled to release resources.

	// Initialize the application and get the container with all dependencies.
	container := dependency.CreateApplication(ctx)
	if container == nil {
		// If the application failed to initialize, log the error and exit with a non-zero status.
		logging.Logger(failedInitialization)
		os.Exit(1)
	}

	// Launch the server in a separate goroutine.
	go func() {
		// Launch the server using the delivery factory from the container.
		container.DeliveryFactory.LaunchServer(ctx, container)
	}()

	// Set up a channel to listen for OS signals (e.g., SIGINT, SIGTERM).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	<-quit

	// Perform a graceful shutdown when a signal is received.
	applicationModel.GracefulShutdown(ctx, container)
}
