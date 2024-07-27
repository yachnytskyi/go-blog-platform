package dependency

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// NewApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func NewApplication(ctx context.Context) model.Container {
	config := factory.NewConfig(constants.Config)
	logger := factory.NewLogger(ctx, config)

	// Create repositories
	repositoryFactory := factory.NewRepositoryFactory(ctx, config, logger)
	repository := repositoryFactory.NewRepository(ctx)
	userRepository := repositoryFactory.NewUserRepository(repository)
	postRepository := repositoryFactory.NewPostRepository(repository)

	// Create use cases
	usecaseFactory := factory.NewUseCaseFactory(ctx, config, logger, repositoryFactory)
	userUseCase := usecaseFactory.NewUserUseCase(userRepository)
	postUseCase := usecaseFactory.NewPostUseCase(postRepository)

	// Create controllers
	deliveryFactory := factory.NewDeliveryFactory(ctx, config, logger, repositoryFactory)
	userController := deliveryFactory.NewUserController(userUseCase)
	postController := deliveryFactory.NewPostController(userUseCase, postUseCase)

	container := model.NewContainer(logger, repositoryFactory, usecaseFactory, deliveryFactory)
	serverRouters := model.NewServerRouters(
		userUseCase,
		deliveryFactory.NewUserRouter(userController),
		deliveryFactory.NewPostRouter(postController),
	)

	container.Delivery.NewDelivery(serverRouters)
	return container
}
