package domain

import (
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/usecase"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userUseCase "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
)

// UseCaseV1 is a struct responsible for creating version 1 instances of the use cases.
// This struct serves as a concrete implementation of the UseCase interface, focusing on the V1 version of use cases.
type UseCaseV1 struct {
}

func NewUseCaseV1() *UseCaseV1 {
	return &UseCaseV1{}
}

// NewUserUseCase creates and returns a new UserUseCase instance using the provided repository.
//
// Parameters:
// - repository (any): The repository instance to be used by the use case. It must implement the user.UserRepository interface.
//
// Returns:
// - user.UserUseCase: The newly created UserUseCase instance.
//
// This method type asserts the provided repository to a user.UserRepository and then creates a new UserUseCaseV1 instance
// with it, ensuring that the use case has access to the necessary repository methods.
func (useCaseV1 UseCaseV1) NewUserUseCase(repository any) user.UserUseCase {
	// Type assert the repository to user.UserRepository.
	userRepository := repository.(user.UserRepository)

	// Create and return a new UserUseCaseV1 instance.
	return userUseCase.NewUserUseCaseV1(userRepository)
}

// NewPostUseCase creates and returns a new PostUseCase instance using the provided repository.
//
// Parameters:
// - repository (any): The repository instance to be used by the use case. It must implement the post.PostRepository interface.
//
// Returns:
// - post.PostUseCase: The newly created PostUseCase instance.
//
// This method type asserts the provided repository to a post.PostRepository and then creates a new PostUseCaseV1 instance
// with it, ensuring that the use case has access to the necessary repository methods.
func (useCaseV1 UseCaseV1) NewPostUseCase(repository any) post.PostUseCase {
	// Type assert the repository to post.PostRepository.
	postRepository := repository.(post.PostRepository)

	// Create and return a new PostUseCaseV1 instance.
	return postUseCase.NewPostUseCaseV1(postRepository)
}
