package dependency

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/application"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location = "pkg.dependendency.CreateApplication"
)

type Container struct {
}

func CreateApplication(ctx context.Context) Container {
	loadConfig, loadConfigError := config.LoadConfig(".")
	if loadConfigError != nil {
		loadConfigInternalError := domainError.NewInternalError(location+".LoadConfig", loadConfigError.Error())
		logging.Logger(loadConfigInternalError)
		application.GracefulShutdown()
	}

	repositoryFactory := repository.InjectRepository(loadConfig)
	db := repositoryFactory.NewRepository(ctx)
	userRepository := repositoryFactory.NewUserRepository(db)
	postRepository := repositoryFactory.NewPostRepository(db)
	fmt.Println(userRepository, postRepository)

	domainFactory := domain.InjectDomain(loadConfig, repositoryFactory)
	userDomain := domainFactory.NewUserRepository(userRepository)
	postDomain := domainFactory.NewPostRepository(postRepository)
	fmt.Println(userDomain, postDomain)
	return Container{}
}
