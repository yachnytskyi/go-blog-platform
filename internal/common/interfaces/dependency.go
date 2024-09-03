package interfaces

import (
	"context"

	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
)

// Logger is an interface that defines methods for logging at different levels.
type Logger interface {
	Trace(data error)
	Debug(data error)
	Info(data error)
	Warn(data error)
	Error(data error)
	Fatal(data error)
	Panic(data error)
}

// Email is an interface that defines methods for sending emails.
type Email interface {
	SendEmail(config *config.ApplicationConfig, logger Logger, location string, data any, emailData EmailData) error
}

// Repository is an interface that defines methods for creating and managing repository instances.
type Repository interface {
	CreateRepository(ctx context.Context) any
	NewRepository(createRepository any, repository any) any
	Close
}

// UseCase is an interface that defines methods for creating use case instances.
type UseCase interface {
	NewUseCase(email Email, repository any) any
}

// Delivery is an interface that defines methods for creating delivery components and managing the server.
type Delivery interface {
	CreateDelivery(serverRouters ServerRouters)
	NewController(userUseCase any, usecase any) any
	NewRouter(router any) Router
	LaunchServer(ctx context.Context, repository Repository)
	Close
}

// Close is an interface that defines a method for closing resources or services.
type Close interface {
	Close(ctx context.Context) // Closes resources or services.
}
