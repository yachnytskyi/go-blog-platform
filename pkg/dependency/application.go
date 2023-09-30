package dependency

import (
	"context"
	"fmt"

	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain"
	container "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"

	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
)

func CreateApplication(ctx context.Context) *container.Container {
	container := container.Container{}
	loadConfig := commonUtility.LoadConfig()

	// Repositories.
	repository.InjectRepository(loadConfig, &container)
	db := container.RepositoryFactory.NewRepository(ctx)
	userRepository := container.RepositoryFactory.NewUserRepository(db)
	postRepository := container.RepositoryFactory.NewPostRepository(db)
	fmt.Println(userRepository, postRepository)

	// Domains.
	domain.InjectDomain(loadConfig, &container)
	userDomain := container.DomainFactory.NewUserRepository(userRepository)
	postDomain := container.DomainFactory.NewPostRepository(postRepository)
	fmt.Println(container)
	fmt.Println(userDomain, postDomain)
	return &container
}
