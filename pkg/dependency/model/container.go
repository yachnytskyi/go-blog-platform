package model

import (
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
)

// Container holds the factory interfaces required to initialize and manage dependencies.
type Container struct {
	Logger     interfaces.Logger     // Interface for creating a logger instance.
	Repository interfaces.Repository // Interface for creating repository instances.
	Delivery   interfaces.Delivery   // Interface for creating delivery components and initializing the server.
} // Add other dependencies as needed.

func NewContainer(logger interfaces.Logger, repository interfaces.Repository, delivery interfaces.Delivery) Container {
	return Container{
		Logger:     logger,
		Repository: repository,
		Delivery:   delivery,
	} // Add other dependencies as needed.
}
