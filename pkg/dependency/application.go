package dependency

import (
	"context"

	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

// CreateApplication initializes the application by setting up the container,
// injecting dependencies, and configuring the server.
func CreateApplication(ctx context.Context) *applicationModel.Container {
	container := &applicationModel.Container{}
	serverRouters := applicationModel.ServerRouters{}

	factory.InjectRepository(ctx, container)
	repository := container.Repository.NewRepository(ctx)
	userRepository := container.Repository.NewUserRepository(repository)
	postRepository := container.Repository.NewPostRepository(repository)

	factory.InjectUseCase(ctx, container)
	userUseCase := container.UseCase.NewUserUseCase(userRepository)
	postUseCase := container.UseCase.NewPostUseCase(postRepository)
	serverRouters.UserUseCase = userUseCase

	factory.InjectDelivery(ctx, container)
	userController := container.Delivery.NewUserController(userUseCase)
	postController := container.Delivery.NewPostController(userUseCase, postUseCase)

	serverRouters.UserRouter = container.Delivery.NewUserRouter(userController)
	serverRouters.PostRouter = container.Delivery.NewPostRouter(postController)

	container.Delivery.InitializeServer(serverRouters)
	return container
}
