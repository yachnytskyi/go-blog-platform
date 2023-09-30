package dependency

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain"
	container "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/application"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location = "pkg.dependendency.CreateApplication"
)

func CreateApplication(ctx context.Context) *container.Container {
	container := container.Container{}
	loadConfig, loadConfigError := config.LoadConfig(config.ConfigPath)
	if loadConfigError != nil {
		loadConfigInternalError := domainError.NewInternalError(location+".LoadConfig", loadConfigError.Error())
		logging.Logger(loadConfigInternalError)
		application.GracefulShutdown(&container)
	}

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
