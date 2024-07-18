package model

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
)

// Container holds the factory interfaces required to initialize and manage dependencies.
type Container struct {
	Repository Repository // Interface for creating repository instances.
	UseCase    UseCase    // Interface for creating use cases.
	Delivery   Delivery   // Interface for creating delivery components and initializing the server.
}

func NewContainer(repository Repository, useCase UseCase, delivery Delivery) *Container {
	return &Container{
		Repository: repository,
		UseCase:    useCase,
		Delivery:   delivery,
	}
}

type ServerRouters struct {
	UserUseCase user.UserUseCase // UserUseCase handles user-related logic and operations.
	UserRouter  user.UserRouter
	PostRouter  post.PostRouter
}

type Repository interface {
	NewRepository(ctx context.Context) any
	CloseRepository(ctx context.Context)
	NewUserRepository(db any) user.UserRepository
	NewPostRepository(db any) post.PostRepository
}

type UseCase interface {
	NewUserUseCase(repository any) user.UserUseCase
	NewPostUseCase(repository any) post.PostUseCase
}

type Delivery interface {
	InitializeServer(serverConfig ServerRouters)
	LaunchServer(ctx context.Context, container *Container)
	CloseServer(ctx context.Context)
	NewUserController(useCase any) user.UserController
	NewPostController(userUseCase, postUseCase any) post.PostController
	NewUserRouter(controller any) user.UserRouter
	NewPostRouter(controller any) post.PostRouter
}
