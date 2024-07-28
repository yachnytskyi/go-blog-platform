package interfaces

import (
	"context"

	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
)

// Config is an interface that defines a method for retrieving the application's configuration.
type Config interface {
	GetConfig() *config.ApplicationConfig
}

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

// Repository is an interface that defines methods for creating and managing repository instances.
type Repository interface {
	CreateRepository(ctx context.Context) any
	NewRepository(createRepository any, repository any) any
	Close
}

// UseCase is an interface that defines methods for creating use case instances.
type UseCase interface {
	NewUseCase(repository any) any
}

// Delivery is an interface that defines methods for creating delivery components and managing the server.
type Delivery interface {
	CreateDelivery(serverRouters ServerRouters)
	NewController(userUseCase UserUseCase, usecase any) any
	NewRouter(router any) any
	LaunchServer(ctx context.Context, repository Repository)
	Close
}

// Close is an interface that defines a method for closing resources or services.
type Close interface {
	Close(ctx context.Context) // Closes resources or services.
}

// ServerRouters holds the routers for different modules of the application.
type ServerRouters struct {
	UserUseCase UserUseCase // UserUseCase handles user-related logic and operations.
	UserRouter  UserRouter
	PostRouter  PostRouter
} // Add other routers as needed.

func NewServerRouters(userUseCase UserUseCase, userRouter UserRouter, postRouter PostRouter) ServerRouters {
	return ServerRouters{
		UserUseCase: userUseCase,
		UserRouter:  userRouter,
		PostRouter:  postRouter,
	} // Add other routers as needed.
}
