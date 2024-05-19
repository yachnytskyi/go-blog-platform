package usecase

import (
	"context"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "internal.post.domain.usecase."
)

type PostUseCase struct {
	postRepository post.PostRepository
}

func NewPostUseCase(postRepository post.PostRepository) post.PostUseCase {
	return &PostUseCase{postRepository: postRepository}
}

func (postUseCase *PostUseCase) GetAllPosts(ctx context.Context, page int, limit int) (*postModel.Posts, error) {
	fetchedPosts, err := postUseCase.postRepository.GetAllPosts(ctx, page, limit)
	return fetchedPosts, err
}

func (postUseCase *PostUseCase) GetPostById(ctx context.Context, postID string) (*postModel.Post, error) {
	fetchedPost, err := postUseCase.postRepository.GetPostById(ctx, postID)
	return fetchedPost, err
}

func (postUseCase *PostUseCase) CreatePost(ctx context.Context, post *postModel.PostCreate) (*postModel.Post, error) {
	createdPost, err := postUseCase.postRepository.CreatePost(ctx, post)
	return createdPost, err
}

func (postUseCase *PostUseCase) UpdatePostById(ctx context.Context, postID string, post *postModel.PostUpdate, currentUserID string) (*postModel.Post, error) {
	fetchedPost, err := postUseCase.GetPostById(ctx, postID)

	if err != nil {
		return nil, err
	}

	userID := fetchedPost.UserID
	if validator.AreStringsNotEqual(currentUserID, userID) {
		return nil, domainError.NewAuthorizationError(location+"UpdatePostById.AreStringsNotEqual", constants.AuthorizationErrorNotification)
	}

	updatedPost, err := postUseCase.postRepository.UpdatePostById(ctx, postID, post)
	return updatedPost, err
}

func (postUseCase *PostUseCase) DeletePostByID(ctx context.Context, postID string, currentUserID string) error {
	fetchedPost, err := postUseCase.GetPostById(ctx, postID)

	if err != nil {
		return err
	}

	userID := fetchedPost.UserID
	if validator.AreStringsNotEqual(currentUserID, userID) {
		return domainError.NewAuthorizationError(location+"DeletePostByID.AreStringsNotEqual", constants.AuthorizationErrorNotification)
	}

	deletedPost := postUseCase.postRepository.DeletePostByID(ctx, postID)
	return deletedPost
}
