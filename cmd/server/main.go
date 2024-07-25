package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	dependency "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

func main() {
	// Create a context with a timeout for the entire application lifecycle.
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimer)
	defer cancel() // Ensure the context is canceled to release resources.

	// Initialize the application and get the container with all dependencies.
	container := dependency.CreateApplication(ctx)

	// Launch the server in a separate goroutine.
	go func() {
		// Launch the server using the delivery from the container.
		container.Delivery.LaunchServer(ctx, container.Repository)
	}()

	// Set up a channel to listen for OS signals (e.g., SIGINT, SIGTERM).
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Block until a signal is received.
	<-quit

	// Perform a graceful shutdown when a signal is received.
	model.GracefulShutdown(ctx, container.Logger, container.Repository, container.Delivery)
}
