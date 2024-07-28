package domain

import (
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
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

func (useCaseV1 UseCaseV1) NewUserUseCase(repository any) user.UserUseCase {
	userRepository := repository.(user.UserRepository)
	return userUseCase.NewUserUseCaseV1(useCaseV1.Config, useCaseV1.Logger, userRepository)
}

func (useCaseV1 UseCaseV1) NewPostUseCase(repository any) post.PostUseCase {
	postRepository := repository.(post.PostRepository)
	return postUseCase.NewPostUseCaseV1(useCaseV1.Logger, postRepository)
}
