package dependency

import (
	"context"

	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// CreateApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
//
// Returns:
// - *applicationModel.Container: The initialized container with all dependencies and server configurations.
func CreateApplication(ctx context.Context) *applicationModel.Container {
	// Initialize a new container for holding the factories.
	container := &applicationModel.Container{}

	// Initialize a structure to hold the server routers.
	serverRouters := applicationModel.ServerRouters{}

	// Inject repository dependencies into the container.
	factory.InjectRepository(ctx, container)

	// Inject repository dependencies into the container.
	repository := container.Repository.NewRepository(ctx)

	// Create specific repositories using the repository instance.
	userRepository := container.Repository.NewUserRepository(repository)
	postRepository := container.Repository.NewPostRepository(repository)

	// Inject use case dependencies into the container.
	factory.InjectUseCase(ctx, container)

	// Create the use cases using the repositories.
	userUseCase := container.UseCase.NewUserUseCase(userRepository)
	postUseCase := container.UseCase.NewPostUseCase(postRepository)
	serverRouters.UserUseCase = userUseCase

	// Inject delivery dependencies into the container.
	factory.InjectDelivery(ctx, container)

	// Create the controllers using the use cases.
	userController := container.Delivery.NewUserController(userUseCase)
	postController := container.Delivery.NewPostController(userUseCase, postUseCase)

	// Initialize the routers using the controllers.
	serverRouters.UserRouter = container.Delivery.NewUserRouter(userController)
	serverRouters.PostRouter = container.Delivery.NewPostRouter(postController)

	// Initialize the server with the configured routers.
	container.Delivery.InitializeServer(serverRouters)

	// Return the initialized container.
	return container
}
