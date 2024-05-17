package dependency

import (
	"context"

	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// CreateApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func CreateApplication(ctx context.Context) *applicationModel.Container {
	// Initialize a new container for holding factory interfaces.
	// Initialize a structure to hold the server routers.
	container := &applicationModel.Container{}
	serverRouters := applicationModel.ServerRouters{}

	// Inject repository dependencies into the container.
	factory.InjectRepository(ctx, container)

	// Create a new repository instance.
	db := container.RepositoryFactory.NewRepository(ctx)
	// If the database connection fails, return nil to indicate the failure.
	if db == nil {
		return nil
	}

	// Create specific repositories.
	userRepository := container.RepositoryFactory.NewUserRepository(db)
	postRepository := container.RepositoryFactory.NewPostRepository(db)

	// Inject domain dependencies into the container.
	factory.InjectDomain(ctx, container)

	// Create the use cases for using the repositories.
	userUseCase := container.DomainFactory.NewUserUseCase(userRepository)
	postUseCase := container.DomainFactory.NewPostUseCase(postRepository)

	// Inject delivery dependencies into the container.
	factory.InjectDelivery(ctx, container)

	// Create the controllers using the use cases.
	userController := container.DeliveryFactory.NewUserController(userUseCase)
	postController := container.DeliveryFactory.NewPostController(postUseCase)

	// Initialize the routers using the controllers.
	serverRouters.UserRouter = container.DeliveryFactory.NewUserRouter(userController)
	serverRouters.PostRouter = container.DeliveryFactory.NewPostRouter(postController)

	// Initialize the server with the configured routers.
	container.DeliveryFactory.InitializeServer(serverRouters)

	// Return the initialized container.
	return container
}
