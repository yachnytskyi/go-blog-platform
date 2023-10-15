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
)

func init() {
	config.LoadConfig()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), constants.DefaultContextTimer)
	defer cancel()
	container := dependency.CreateApplication(ctx)
	go func() {
		container.DeliveryFactory.LaunchServer(ctx, container)
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Perform Graceful Shutdown
	applicationModel.GracefulShutdown(ctx, container)
}
