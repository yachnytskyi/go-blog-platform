package domain

import (
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
)

type UseCaseV1 struct {
	Logger applicationModel.Logger
}

func NewUseCaseV1(logger applicationModel.Logger) *UseCaseV1 {
	return &UseCaseV1{Logger: logger}
}

func (useCaseV1 UseCaseV1) NewUserUseCase(repository any) user.UserUseCase {
	userRepository := repository.(user.UserRepository)
	return userUseCase.NewUserUseCaseV1(useCaseV1.Logger, userRepository)
}

func (useCaseV1 UseCaseV1) NewPostUseCase(repository any) post.PostUseCase {
	postRepository := repository.(post.PostRepository)
	return postUseCase.NewPostUseCaseV1(useCaseV1.Logger, postRepository)
}
