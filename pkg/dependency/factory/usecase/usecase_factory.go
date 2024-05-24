package domain

import (
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
)

// UseCaseFactory is a factory for creating use case instances.
type UseCaseFactory struct {
}

// NewUserUseCase creates and returns a new UserUseCase instance using the provided repository.
func (useCaseFactory UseCaseFactory) NewUserUseCase(repository any) user.UserUseCase {
	// Type assert the repository to user.UserRepository.
	userRepository := repository.(user.UserRepository)

	// Create and return a new UserUseCase instance.
	return userUseCase.NewUserUseCase(userRepository)
}

// NewPostUseCase creates and returns a new PostUseCase instance using the provided repository.
func (useCaseFactory UseCaseFactory) NewPostUseCase(repository any) post.PostUseCase {
	// Type assert the repository to post.PostRepository.
	postRepository := repository.(post.PostRepository)

	// Create and return a new PostUseCase instance.
	return postUseCase.NewPostUseCase(postRepository)
}
