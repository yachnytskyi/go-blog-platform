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

type PostUseCaseV1 struct {
	postRepository post.PostRepository
}

func NewPostUseCaseV1(postRepository post.PostRepository) post.PostUseCase {
	return &PostUseCaseV1{postRepository: postRepository}
}

func (postUseCaseV1 *PostUseCaseV1) GetAllPosts(ctx context.Context, page int, limit int) (*postModel.Posts, error) {
	fetchedPosts, err := postUseCaseV1.postRepository.GetAllPosts(ctx, page, limit)
	return fetchedPosts, err
}

func (postUseCaseV1 *PostUseCaseV1) GetPostById(ctx context.Context, postID string) (*postModel.Post, error) {
	fetchedPost, err := postUseCaseV1.postRepository.GetPostById(ctx, postID)
	return fetchedPost, err
}

func (postUseCaseV1 *PostUseCaseV1) CreatePost(ctx context.Context, post *postModel.PostCreate) (*postModel.Post, error) {
	createdPost, err := postUseCaseV1.postRepository.CreatePost(ctx, post)
	return createdPost, err
}

func (postUseCaseV1 *PostUseCaseV1) UpdatePostById(ctx context.Context, postID string, post *postModel.PostUpdate, currentUserID string) (*postModel.Post, error) {
	fetchedPost, err := postUseCaseV1.GetPostById(ctx, postID)

	if err != nil {
		return nil, err
	}

	userID := fetchedPost.UserID
	if validator.AreStringsNotEqual(currentUserID, userID) {
		return nil, domainError.NewAuthorizationError(location+"UpdatePostById.AreStringsNotEqual", constants.AuthorizationErrorNotification)
	}

	updatedPost, err := postUseCaseV1.postRepository.UpdatePostById(ctx, postID, post)
	return updatedPost, err
}

func (postUseCaseV1 *PostUseCaseV1) DeletePostByID(ctx context.Context, postID string, currentUserID string) error {
	fetchedPost, err := postUseCaseV1.GetPostById(ctx, postID)

	if err != nil {
		return err
	}

	userID := fetchedPost.UserID
	if validator.AreStringsNotEqual(currentUserID, userID) {
		return domainError.NewAuthorizationError(location+"DeletePostByID.AreStringsNotEqual", constants.AuthorizationErrorNotification)
	}

	deletedPost := postUseCaseV1.postRepository.DeletePostByID(ctx, postID)
	return deletedPost
}
