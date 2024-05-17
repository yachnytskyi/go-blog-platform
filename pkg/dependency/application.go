package dependency

import (
	"context"

	factory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

func CreateApplication(ctx context.Context) *applicationModel.Container {
	container := &applicationModel.Container{}
	serverRouters := applicationModel.ServerRouters{}

	// Repositories.
	factory.InjectRepository(ctx, container)
	db := container.RepositoryFactory.NewRepository(ctx)
	userRepository := container.RepositoryFactory.NewUserRepository(db)
	postRepository := container.RepositoryFactory.NewPostRepository(db)

	// Domains.
	factory.InjectDomain(ctx, container)
	userUseCase := container.DomainFactory.NewUserUseCase(userRepository)
	postDomain := container.DomainFactory.NewPostUseCase(postRepository)

	// Deliveries.
	factory.InjectDelivery(ctx, container)
	userController := container.DeliveryFactory.NewUserController(userUseCase)
	postController := container.DeliveryFactory.NewPostController(postDomain)
	serverRouters.UserRouter = container.DeliveryFactory.NewUserRouter(userController)
	serverRouters.PostRouter = container.DeliveryFactory.NewPostRouter(postController)
	container.DeliveryFactory.InitializeServer(serverRouters)
	return container
}
