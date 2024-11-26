package usecase

import (
	"context"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/interfaces"
	model "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location = "internal.post.domain.usecase."
)

type PostUseCase struct {
	Logger         interfaces.Logger
	PostRepository interfaces.PostRepository
}

func NewPostUseCase(logger interfaces.Logger, postRepository interfaces.PostRepository) PostUseCase {
	return PostUseCase{
		Logger:         logger,
		PostRepository: postRepository,
	}
}

func (postUseCase PostUseCase) GetAllPosts(ctx context.Context, page int, limit int) (*model.Posts, error) {
	fetchedPosts, err := postUseCase.PostRepository.GetAllPosts(ctx, page, limit)
	return fetchedPosts, err
}

func (postUseCase PostUseCase) GetPostById(ctx context.Context, postID string) (*model.Post, error) {
	fetchedPost, err := postUseCase.PostRepository.GetPostById(ctx, postID)
	return fetchedPost, err
}

func (postUseCase PostUseCase) CreatePost(ctx context.Context, post *model.PostCreate) (*model.Post, error) {
	createdPost, err := postUseCase.PostRepository.CreatePost(ctx, post)
	return createdPost, err
}

func (postUseCase PostUseCase) UpdatePostById(ctx context.Context, postID string, post *model.PostUpdate, currentUserID string) (*model.Post, error) {
	fetchedPost, err := postUseCase.GetPostById(ctx, postID)

	if err != nil {
		return nil, err
	}

	userID := fetchedPost.UserID
	if currentUserID != userID {
		return nil, domain.NewAuthorizationError(location, constants.AuthorizationErrorNotification)
	}

	updatedPost, err := postUseCase.PostRepository.UpdatePostById(ctx, postID, post)
	return updatedPost, err
}

func (postUseCase PostUseCase) DeletePostByID(ctx context.Context, postID string, currentUserID string) error {
	fetchedPost, err := postUseCase.GetPostById(ctx, postID)

	if err != nil {
		return err
	}

	userID := fetchedPost.UserID
	if currentUserID != userID {
		return domain.NewAuthorizationError(location, constants.AuthorizationErrorNotification)
	}

	deletedPost := postUseCase.PostRepository.DeletePostByID(ctx, postID)
	return deletedPost
}
