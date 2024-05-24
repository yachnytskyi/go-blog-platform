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

	// Create a new repository instance using the repository factory in the container.
	db := container.RepositoryFactory.NewRepository(ctx)

	// Create specific repositories using the repository instance.
	userRepository := container.RepositoryFactory.NewUserRepository(db)
	postRepository := container.RepositoryFactory.NewPostRepository(db)

	// Inject use case dependencies into the container.
	factory.InjectUseCase(ctx, container)

	// Create the use cases using the repositories.
	userUseCase := container.UseCaseFactory.NewUserUseCase(userRepository)
	postUseCase := container.UseCaseFactory.NewPostUseCase(postRepository)

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
