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

func NewContainer(repositoryFactory RepositoryFactory, domainFactory DomainFactory, deliveryFactory DeliveryFactory) *Container {
	return &Container{
		RepositoryFactory: repositoryFactory,
		DomainFactory:     domainFactory,
		DeliveryFactory:   deliveryFactory,
	}
}

// Define a DatabaseFactory interface to create different database instances.
type RepositoryFactory interface {
	NewRepository(ctx context.Context) interface{}
	CloseRepository()
	NewUserRepository(db interface{}) user.UserRepository
	NewPostRepository(db interface{}) post.PostRepository
}

// Define a DomainFactory interface to create different domain instances.
type DomainFactory interface {
	NewUserUseCase(repository interface{}) user.UserUseCase
	NewPostUseCase(repository interface{}) post.PostUseCase
}

// Define a DatabaseFactory interface to create different database instances.
type DeliveryFactory interface {
	NewUserController(domain interface{}) user.UserController
	NewPostController(domain interface{}) post.PostController

	CloseServer()
}
