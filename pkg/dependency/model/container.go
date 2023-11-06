package model

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
)

type Container struct {
	RepositoryFactory RepositoryFactory
	DomainFactory     DomainFactory
	DeliveryFactory   DeliveryFactory
}

type ServerRouters struct {
	UserUseCase user.UserUseCase
	UserRouter  user.UserRouter
	PostRouter  post.PostRouter
}

func NewContainer(repositoryFactory RepositoryFactory, domainFactory DomainFactory, deliveryFactory DeliveryFactory) *Container {
	return &Container{
		RepositoryFactory: repositoryFactory,
		DomainFactory:     domainFactory,
		DeliveryFactory:   deliveryFactory,
	}
}

// Define a DatabaseFactory interface to create different database instances.
type RepositoryFactory interface {
	NewRepository(ctx context.Context) any
	CloseRepository(ctx context.Context)
	NewUserRepository(db any) user.UserRepository
	NewPostRepository(db any) post.PostRepository
}

// Define a DomainFactory interface to create different domain instances.
type DomainFactory interface {
	NewUserUseCase(repository any) user.UserUseCase
	NewPostUseCase(repository any) post.PostUseCase
}

// Define a DatabaseFactory interface to create different database instances.
type DeliveryFactory interface {
	InitializeServer(serverConfig ServerRouters)
	LaunchServer(ctx context.Context, container *Container)
	CloseServer(ctx context.Context)
	NewUserController(domain any) user.UserController
	NewPostController(domain any) post.PostController
	NewUserRouter(controller any) user.UserRouter
	NewPostRouter(controller any) post.PostRouter
}
