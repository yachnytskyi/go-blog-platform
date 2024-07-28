package model

import (
	"context"

	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
)

// Container holds the factory interfaces required to initialize and manage dependencies.
type Container struct {
	Logger     interfaces.Logger // Interface for creating a logger instance.
	Repository Repository        // Interface for creating repository instances.
	Delivery   Delivery          // Interface for creating delivery components and initializing the server.
} // Add other dependencies as needed

// ServerRouters holds the routers for different modules of the application.
type ServerRouters struct {
	UserUseCase interfaces.UserUseCase // UserUseCase handles user-related logic and operations.
	UserRouter  interfaces.UserRouter
	PostRouter  interfaces.PostRouter
} // Add other routers as needed

func NewContainer(logger interfaces.Logger, repository Repository, delivery Delivery) Container {
	return Container{
		Logger:     logger,
		Repository: repository,
		Delivery:   delivery,
	} // Add other dependencies as needed
}

func NewServerRouters(userUseCase interfaces.UserUseCase, userRouter interfaces.UserRouter, postRouter interfaces.PostRouter) ServerRouters {
	return ServerRouters{
		UserUseCase: userUseCase,
		UserRouter:  userRouter,
		PostRouter:  postRouter,
	} // Add other routers as needed
}

// Repository is an interface that defines methods for creating and managing repository instances.
type Repository interface {
	NewRepository(ctx context.Context) any
	Closer
	NewUserRepository(repository any) interfaces.UserRepository
	NewPostRepository(repository any) interfaces.PostRepository
}

// UseCase is an interface that defines methods for creating use case instances.
type UseCase interface {
	NewUserUseCase(repository any) interfaces.UserUseCase
	NewPostUseCase(repository any) interfaces.PostUseCase
}

// Delivery is an interface that defines methods for creating delivery components and managing the server.
type Delivery interface {
	NewDelivery(serverRouters ServerRouters)
	LaunchServer(ctx context.Context, repository Repository)
	Closer
	NewUserController(useCase any) interfaces.UserController
	NewPostController(userUseCase, postUseCase any) interfaces.PostController
	NewUserRouter(controller any) interfaces.UserRouter
	NewPostRouter(controller any) interfaces.PostRouter
}

// Closer is an interface that defines a method for closing resources or services.
type Closer interface {
	Close(ctx context.Context) // Closes resources or services.
}
