package service

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"
)

type PostService struct {
	postRepository post.Repository
}

func NewPostService(postRepository post.Repository) post.Service {
	return &PostService{postRepository: postRepository}
}

func (postService *PostService) GetPostById(ctx context.Context, postID string) (*models.PostDB, error) {
	fetchedPost, err := postService.postRepository.GetPostById(ctx, postID)

	return fetchedPost, err
}

func (postService *PostService) GetAllPosts(ctx context.Context, page int, limit int) ([]*models.PostDB, error) {
	fetchedPosts, err := postService.postRepository.GetAllPosts(ctx, page, limit)

	return fetchedPosts, err
}

func (postService *PostService) CreatePost(ctx context.Context, post *models.PostCreate) (*models.PostDB, error) {
	createdPost, err := postService.postRepository.CreatePost(ctx, post)

	return createdPost, err
}

func (postService *PostService) UpdatePostById(ctx context.Context, postID string, post *models.PostUpdate, currentUserID string) (*models.PostDB, error) {
	fetchedPost, err := postService.GetPostById(ctx, postID)

	if err != nil {
		return nil, err
	}

	userID := fetchedPost.UserID

	if err := utils.IsUserOwner(currentUserID, userID); err != nil {
		return nil, err
	}

	updatedPost, err := postService.postRepository.UpdatePostById(ctx, postID, post)

	return updatedPost, err
}

func (postService *PostService) DeletePostByID(ctx context.Context, postID string, currentUserID string) error {
	fetchedPost, err := postService.GetPostById(ctx, postID)

	if err != nil {
		return err
	}

	userID := fetchedPost.UserID

	if err := utils.IsUserOwner(currentUserID, userID); err != nil {
		return err
	}

	deletedPost := postService.postRepository.DeletePostByID(ctx, postID)

	return deletedPost
}
