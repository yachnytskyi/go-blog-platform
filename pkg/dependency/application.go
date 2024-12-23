package dependency

import (
	"context"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
)

// NewApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func NewApplication(ctx context.Context) model.Container {
	config := factory.NewConfig(constants.Config)
	logger := factory.NewLogger(config)
	email := factory.NewEmail(config, logger)

	// Create repository factory and repositories, then assert their types.
	repository := factory.NewRepositoryFactory(config, logger)
	createRepository := repository.CreateRepository(ctx)
	userRepository := repository.NewRepository(createRepository, (*interfaces.UserRepository)(nil)).(interfaces.UserRepository)
	postRepository := repository.NewRepository(createRepository, (*interfaces.PostRepository)(nil)).(interfaces.PostRepository)

	// Create use cases.
	userUseCase := user.NewUserUseCase(config, logger, email, userRepository)
	postUseCase := post.NewPostUseCase(logger, postRepository)

	// Create delivery factory and controllers.
	delivery := factory.NewDeliveryFactory(ctx, config, logger, repository)
	healthController := delivery.NewHealthCheckController(repository)
	userController := delivery.NewController(userUseCase)
	postController := delivery.NewController(postUseCase)

	// Create routers.
	serverRouters := interfaces.NewServerRouters(
		delivery.NewHealthRouter(healthController, repository),
		delivery.NewRouter(userController),
		delivery.NewRouter(postController),
		// Add other routers as needed.
	)

	delivery.CreateDelivery(serverRouters)
	repository.HealthCheck(delivery)
	container := model.NewContainer(logger, repository, delivery)
	return container
}
