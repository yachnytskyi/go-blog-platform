package dependency

import (
	"context"
	"fmt"

	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/delivery"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

func CreateApplication(ctx context.Context) *applicationModel.Container {
	container := applicationModel.Container{}

	// Repositories.
	repository.InjectRepository(&container)
	db := container.RepositoryFactory.NewRepository(ctx)
	userRepository := container.RepositoryFactory.NewUserRepository(db)
	postRepository := container.RepositoryFactory.NewPostRepository(db)
	fmt.Println(userRepository, postRepository)

	// Domains.
	domain.InjectDomain(&container)
	userDomain := container.DomainFactory.NewUserUseCase(userRepository)
	postDomain := container.DomainFactory.NewPostUseCase(postRepository)
	fmt.Println(userDomain, postDomain)

	// Deliveries.
	delivery.InjectDelivery(&container)
	userDelivery := container.DeliveryFactory.NewUserController(userDomain)
	postDelivery := container.DeliveryFactory.NewPostController(postDomain)
	fmt.Println(container)
	fmt.Println(userDelivery, postDelivery)

	// Routers

	return &container
}
