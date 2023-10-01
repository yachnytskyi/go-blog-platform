package dependency

import (
	"context"
	"fmt"

	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

func CreateApplication(ctx context.Context) *applicationModel.Container {
	container := applicationModel.Container{}
	applicationConfig := applicationModel.ApplicationConfig

	// Repositories.
	repository.InjectRepository(applicationConfig, &container)
	db := container.RepositoryFactory.NewRepository(ctx)
	userRepository := container.RepositoryFactory.NewUserRepository(db)
	postRepository := container.RepositoryFactory.NewPostRepository(db)
	fmt.Println(userRepository, postRepository)

	// Domains.
	domain.InjectDomain(applicationConfig, &container)
	userDomain := container.DomainFactory.NewUserRepository(userRepository)
	postDomain := container.DomainFactory.NewPostRepository(postRepository)
	fmt.Println(container)
	fmt.Println(userDomain, postDomain)
	return &container
}
