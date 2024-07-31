package dependency

import (
	"context"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// NewApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func NewApplication(ctx context.Context) model.Container {
	config := factory.NewConfig(constants.Config)
	logger := factory.NewLogger(config)
	email := factory.NewEmail(config, logger)

	// Create repository factory and repositories
	repositoryFactory := factory.NewRepositoryFactory(config, logger)
	createRepository := repositoryFactory.CreateRepository(ctx)
	userRepository := repositoryFactory.NewRepository(createRepository, (*interfaces.UserRepository)(nil))
	postRepository := repositoryFactory.NewRepository(createRepository, (*interfaces.PostRepository)(nil))

	// Create use case factory and use cases.
	usecaseFactory := factory.NewUseCaseFactory(ctx, config, logger, email, repositoryFactory)
	userUseCase := usecaseFactory.NewUseCase(email, userRepository).(interfaces.UserUseCase)
	postUseCase := usecaseFactory.NewUseCase(email, postRepository)

	// Create delivery factory and controllers.
	deliveryFactory := factory.NewDeliveryFactory(ctx, config, logger, repositoryFactory)
	userController := deliveryFactory.NewController(userUseCase, nil)
	postController := deliveryFactory.NewController(userUseCase, postUseCase)

	// Create routers.
	serverRouters := interfaces.NewServerRouters(
		deliveryFactory.NewRouter(userController),
		deliveryFactory.NewRouter(postController),
		// Add other routers as needed.
	)

	deliveryFactory.CreateDelivery(serverRouters)
	container := model.NewContainer(logger, repositoryFactory, deliveryFactory)
	return container
}
