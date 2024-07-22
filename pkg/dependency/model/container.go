package model

import (
	"context"

	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
)

// Container holds the factory interfaces required to initialize and manage dependencies.
type Container struct {
	Logger     Logger     // Interface for creating a logger instance.
	Repository Repository // Interface for creating repository instances.
	UseCase    UseCase    // Interface for creating use cases.
	Delivery   Delivery   // Interface for creating delivery components and initializing the server.
} // Add other dependencies as needed

// ServerRouters holds the routers for different modules of the application.
type ServerRouters struct {
	UserUseCase user.UserUseCase // UserUseCase handles user-related logic and operations.
	UserRouter  user.UserRouter
	PostRouter  post.PostRouter
} // Add other routers as needed

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
	NewRepository(ctx context.Context) any
	CloseRepository(ctx context.Context)
	NewUserRepository(repository any) user.UserRepository
	NewPostRepository(repository any) post.PostRepository
}

// UseCase is an interface that defines methods for creating use case instances.
type UseCase interface {
	NewUserUseCase(repository any) user.UserUseCase
	NewPostUseCase(repository any) post.PostUseCase
}

// Delivery is an interface that defines methods for creating delivery components and managing the server.
type Delivery interface {
	NewDelivery(serverRouters ServerRouters)
	LaunchServer(ctx context.Context, repository Repository)
	CloseServer(ctx context.Context)
	NewUserController(useCase any) user.UserController
	NewPostController(userUseCase, postUseCase any) post.PostController
	NewUserRouter(controller any) user.UserRouter
	NewPostRouter(controller any) post.PostRouter
}

func NewContainer(logger Logger, repository Repository, useCase UseCase, delivery Delivery) *Container {
	return &Container{
		Logger:     logger,
		Repository: repository,
		UseCase:    useCase,
		Delivery:   delivery,
	} // Add other dependencies as needed
}

func NewServerRouters(userUseCase user.UserUseCase, userRouter user.UserRouter, postRouter post.PostRouter) ServerRouters {
	return ServerRouters{
		UserUseCase: userUseCase,
		UserRouter:  userRouter,
		PostRouter:  postRouter,
	} // Add other routers as needed
}
