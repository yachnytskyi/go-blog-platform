package domain

import (
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
)

// UseCaseFactoryV1 is a factory for creating use case instances.
type UseCaseFactoryV1 struct {
}

// NewUserUseCaseV1 creates and returns a new UserUseCaseV1 instance using the provided repository.
func (useCaseFactoryV1 UseCaseFactoryV1) NewUserUseCase(repository any) user.UserUseCase {
	// Type assert the repository to user.UserRepository.
	userRepository := repository.(user.UserRepository)

	// Create and return a new UserUseCaseV1 instance.
	return userUseCase.NewUserUseCaseV1(userRepository)
}

// NewPostUseCaseV1 creates and returns a new PostUseCaseV1 instance using the provided repository.
func (useCaseFactoryV1 UseCaseFactoryV1) NewPostUseCase(repository any) post.PostUseCase {
	// Type assert the repository to post.PostRepository.
	postRepository := repository.(post.PostRepository)

	// Create and return a new PostUseCaseV1 instance.
	return postUseCase.NewPostUseCaseV1(postRepository)
}
