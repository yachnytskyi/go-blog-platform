package v1

import (
	"fmt"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location = "pkg.dependency.factory.data.usecase.v1."
)

type UseCaseV1 struct {
	Config interfaces.Config
	Logger interfaces.Logger
}

func NewUseCaseV1(config interfaces.Config, logger interfaces.Logger) UseCaseV1 {
	return UseCaseV1{
		Config: config,
		Logger: logger,
	}
}

func (useCaseV1 UseCaseV1) NewUseCase(email interfaces.Email, repository any) any {
	switch repositoryType := repository.(type) {
	case interfaces.UserRepository:
		return user.NewUserUseCaseV1(useCaseV1.Config, useCaseV1.Logger, email, repositoryType)
	case interfaces.PostRepository:
		return post.NewPostUseCaseV1(useCaseV1.Logger, repositoryType)
	default:
		useCaseV1.Logger.Panic(domainError.NewInternalError(location+"NewUseCase.default", fmt.Sprintf(constants.UnsupportedUseCase, repositoryType)))
	}

	return nil
}
