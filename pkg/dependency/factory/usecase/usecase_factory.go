package domain

import (
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
)

type UseCaseV1 struct {
}

func NewUseCaseV1() *UseCaseV1 {
	return &UseCaseV1{}
}

func (useCaseV1 UseCaseV1) NewUserUseCase(repository any) user.UserUseCase {
	userRepository := repository.(user.UserRepository)
	return userUseCase.NewUserUseCaseV1(userRepository)
}

func (useCaseV1 UseCaseV1) NewPostUseCase(repository any) post.PostUseCase {
	postRepository := repository.(post.PostRepository)
	return postUseCase.NewPostUseCaseV1(postRepository)
}
