package dependency

import (
	"context"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/interfaces"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// NewApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func NewApplication(ctx context.Context) model.Container {
	config := factory.NewConfig(constants.Config)
	logger := factory.NewLogger(config)
	email := factory.NewEmail(config, logger)

	// Create repository factory and repositories, then assert their types.
	repositoryFactory := factory.NewRepositoryFactory(config, logger)
	createRepository := repositoryFactory.CreateRepository(ctx)
	userRepository := repositoryFactory.NewRepository(createRepository, (*interfaces.UserRepository)(nil)).(interfaces.UserRepository)
	postRepository := repositoryFactory.NewRepository(createRepository, (*interfaces.PostRepository)(nil)).(interfaces.PostRepository)

	// Create use cases.
	userUseCase := user.NewUserUseCase(config, logger, email, userRepository)
	postUseCase := post.NewPostUseCase(logger, postRepository)

	// Create delivery factory and controllers.
	deliveryFactory := factory.NewDeliveryFactory(ctx, config, logger, repositoryFactory)
	userController := deliveryFactory.NewController(userUseCase)
	postController := deliveryFactory.NewController(postUseCase)

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
