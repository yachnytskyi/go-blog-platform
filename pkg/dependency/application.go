package dependency

import (
	"context"

	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/delivery"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

func CreateApplication(ctx context.Context) *applicationModel.Container {
	container := &applicationModel.Container{}
	serverRouters := applicationModel.ServerRouters{}

	// Repositories.
	repository.InjectRepository(ctx, container)
	db := container.RepositoryFactory.NewRepository(ctx)
	userRepository := container.RepositoryFactory.NewUserRepository(db)
	postRepository := container.RepositoryFactory.NewPostRepository(db)

	// Domains.
	domain.InjectDomain(ctx, container)
	serverRouters.UserUseCase = container.DomainFactory.NewUserUseCase(userRepository)
	postDomain := container.DomainFactory.NewPostUseCase(postRepository)

	// Deliveries.
	delivery.InjectDelivery(ctx, container)
	userController := container.DeliveryFactory.NewUserController(serverRouters.UserUseCase)
	postController := container.DeliveryFactory.NewPostController(postDomain)
	serverRouters.UserRouter = container.DeliveryFactory.NewUserRouter(userController)
	serverRouters.PostRouter = container.DeliveryFactory.NewPostRouter(postController)
	container.DeliveryFactory.InitializeServer(ctx, serverRouters)
	return container
}
