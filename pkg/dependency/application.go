package dependency

import (
	"context"

	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// CreateApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func CreateApplication(ctx context.Context) *applicationModel.Container {
	loggerFactory := factory.NewLoggerFactory(ctx)

	// Create repositories
	repositoryFactory := factory.NewRepositoryFactory(ctx)
	repository := repositoryFactory.NewRepository(ctx, loggerFactory)
	userRepository := repositoryFactory.NewUserRepository(repository)
	postRepository := repositoryFactory.NewPostRepository(repository)

	// Create use cases
	usecaseFactory := factory.NewUseCaseFactory(ctx, repositoryFactory)
	userUseCase := usecaseFactory.NewUserUseCase(userRepository)
	postUseCase := usecaseFactory.NewPostUseCase(postRepository)

	// Create controllers
	deliveryFactory := factory.NewDeliveryFactory(ctx, repositoryFactory)
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
