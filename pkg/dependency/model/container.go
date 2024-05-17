package model

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
)

// Container holds the factory interfaces required to initialize and manage dependencies.
type Container struct {
	RepositoryFactory RepositoryFactory // Interface for creating repository instances.
	DomainFactory     DomainFactory     // Interface for creating domain use cases.
	DeliveryFactory   DeliveryFactory   // Interface for creating delivery components and initializing the server.
}

// ServerRouters holds the routers for different entities.
type ServerRouters struct {
	UserRouter user.UserRouter // Router for user-related endpoints.
	PostRouter post.PostRouter // Router for post-related endpoints.
}

// NewContainer initializes and returns a new Container with the provided factories.
func NewContainer(repositoryFactory RepositoryFactory, domainFactory DomainFactory, deliveryFactory DeliveryFactory) *Container {
	return &Container{
		RepositoryFactory: repositoryFactory,
		DomainFactory:     domainFactory,
		DeliveryFactory:   deliveryFactory,
	}
}

// RepositoryFactory defines methods for creating different repository instances and managing their lifecycle.
type RepositoryFactory interface {
	// NewRepository creates and returns a new repository instance.
	NewRepository(ctx context.Context) any
	// CloseRepository closes the repository instance and releases resources.
	CloseRepository(ctx context.Context)
	// NewUserRepository creates and returns a new UserRepository instance.
	NewUserRepository(db any) user.UserRepository
	// NewPostRepository creates and returns a new PostRepository instance.
	NewPostRepository(db any) post.PostRepository
}

// DomainFactory defines methods for creating domain use cases.
type DomainFactory interface {
	// NewUserUseCase creates and returns a new UserUseCase instance using the provided repository.
	NewUserUseCase(repository any) user.UserUseCase
	// NewPostUseCase creates and returns a new PostUseCase instance using the provided repository.
	NewPostUseCase(repository any) post.PostUseCase
}

// DeliveryFactory defines methods for creating delivery components, initializing the server, and managing its lifecycle.
type DeliveryFactory interface {
	// InitializeServer initializes the server with the provided routers configuration.
	InitializeServer(serverConfig ServerRouters)
	// LaunchServer starts the server using the provided context and container.
	LaunchServer(ctx context.Context, container *Container)
	// CloseServer gracefully shuts down the server using the provided context.
	CloseServer(ctx context.Context)
	// NewUserController creates and returns a new UserController instance using the provided domain use case.
	NewUserController(domain any) user.UserController
	// NewPostController creates and returns a new PostController instance using the provided domain use case.
	NewPostController(domain any) post.PostController
	// NewUserRouter creates and returns a new UserRouter instance using the provided controller.
	NewUserRouter(controller any) user.UserRouter
	// NewPostRouter creates and returns a new PostRouter instance using the provided controller.
	NewPostRouter(controller any) post.PostRouter
}
