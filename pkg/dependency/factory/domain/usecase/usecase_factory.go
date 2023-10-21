package usecase

import (
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
)

type UseCaseFactory struct {
}

func (useCaseFactory UseCaseFactory) NewUserUseCase(repository interface{}) user.UserUseCase {
	userRepopository := repository.(user.UserRepository)
	return userUseCase.NewUserUseCase(userRepopository)
}

func (useCaseFactory UseCaseFactory) NewPostUseCase(repository interface{}) post.PostUseCase {
	postRepository := repository.(post.PostRepository)
	return postUseCase.NewPostUseCase(postRepository)
}
