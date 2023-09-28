package domain

import (
	// "context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	useCaseFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/domain/usecase"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location          = "pkg/dependency/domain/InjectDomain"
	unsupportedDomain = "unsupported domain type: %s"
)

// Define a DatabaseFactory interface to create different database instances.
type DomainFactory interface {
	NewUserRepository(db interface{}) user.UserUseCase
	NewPostRepository(db interface{}) post.PostUseCase
}

func InjectDomain(loadConfig config.Config, repositoryFactory repository.RepositoryFactory) DomainFactory {
	switch loadConfig.Domain {
	case config.UseCase:
		return useCaseFactory.UseCaseFactory{}
	// Add other domain options here as needed.
	default:
		logging.Logger(domainError.NewInternalError(location+".loadConfig.Domain:", fmt.Sprintf(unsupportedDomain, loadConfig.Domain)))
		GracefulShutdownDomain(repositoryFactory)
		return nil
	}
}
