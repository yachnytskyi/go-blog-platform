package model

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
)

type Container struct {
	RepositoryFactory RepositoryFactory
	DomainFactory     DomainFactory
}

func NewContainer(repositoryFactory RepositoryFactory, domainFactory DomainFactory) *Container {
	return &Container{
		RepositoryFactory: repositoryFactory,
		DomainFactory:     domainFactory,
	}
}

// Define a DatabaseFactory interface to create different database instances.
type RepositoryFactory interface {
	NewRepository(ctx context.Context) interface{}
	CloseRepository()
	NewUserRepository(db interface{}) user.UserRepository
	NewPostRepository(db interface{}) post.PostRepository
}

// Define a DatabaseFactory interface to create different database instances.
type DomainFactory interface {
	NewUserRepository(db interface{}) user.UserUseCase
	NewPostRepository(db interface{}) post.PostUseCase
}
