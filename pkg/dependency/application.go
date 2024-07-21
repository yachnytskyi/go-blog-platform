package dependency

import (
	"context"

	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// CreateApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func CreateApplication(ctx context.Context) *applicationModel.Container {
	// Initialize the factories
	loggerFactory := factory.NewLoggerFactory(ctx)
	repositoryFactory := factory.NewRepositoryFactory(ctx)
	usecaseFactory := factory.NewUseCaseFactory(ctx, repositoryFactory)
	deliveryFactory := factory.NewDeliveryFactory(ctx, repositoryFactory)

	// Create repositories
	repository := repositoryFactory.NewRepository(ctx, loggerFactory)
	userRepository := repositoryFactory.NewUserRepository(repository)
	postRepository := repositoryFactory.NewPostRepository(repository)

	// Create use cases
	userUseCase := usecaseFactory.NewUserUseCase(userRepository)
	postUseCase := usecaseFactory.NewPostUseCase(postRepository)

	// Create controllers
	userController := deliveryFactory.NewUserController(userUseCase)
	postController := deliveryFactory.NewPostController(userUseCase, postUseCase)

	container := applicationModel.NewContainer(loggerFactory, repositoryFactory, usecaseFactory, deliveryFactory)
	serverRouters := applicationModel.NewServerRouters(
		userUseCase,
		deliveryFactory.NewUserRouter(userController),
		deliveryFactory.NewPostRouter(postController),
	)

	container.Delivery.NewDelivery(serverRouters)
	return container
}
