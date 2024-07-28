package v1

import (
	"fmt"

	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
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

func (useCaseV1 UseCaseV1) NewUseCase(repository any) any {
	switch repositoryType := repository.(type) {
	case interfaces.UserRepository:
		return user.NewUserUseCaseV1(useCaseV1.Config, useCaseV1.Logger, repositoryType)
	case interfaces.PostRepository:
		return post.NewPostUseCaseV1(useCaseV1.Logger, repositoryType)
	default:
		useCaseV1.Logger.Panic(domainError.NewInternalError(location+"NewUseCase.default", fmt.Sprintf(model.UnsupportedUseCase, repositoryType)))
	}

	return nil
}
