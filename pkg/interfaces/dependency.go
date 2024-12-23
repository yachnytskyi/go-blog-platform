package interfaces

import (
	"context"

	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
)

type Logger interface {
	Trace(data error)
	Debug(data error)
	Info(data error)
	Warn(data error)
	Error(data error)
	Fatal(data error)
	Panic(data error)
}

type Email interface {
	SendEmail(config *config.ApplicationConfig, logger Logger, location string, data any, emailData EmailData) error
}

type Repository interface {
	CreateRepository(ctx context.Context) any
	NewRepository(createRepository any, repository any) any
	HealthCheck(delivery Delivery)
	DatabasePing() bool
	Close
}

type Delivery interface {
	CreateDelivery(serverRouters ServerRouters)
	NewHealthCheckController(repository Repository) any
	NewController(useCase any) any
	NewHealthRouter(router any, repository Repository) Router
	NewRouter(router any) Router
	LaunchServer(ctx context.Context, repository Repository)
	Close
}

// Close is an interface that defines a method for closing resources or services.
type Close interface {
	Close(ctx context.Context) // Closes resources or services.
}
