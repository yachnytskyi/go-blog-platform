package usecase

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postModel "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
)

type PostUseCase struct {
	postRepository post.Repository
}

func NewPostUseCase(postRepository post.Repository) post.UseCase {
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

	err = utility.IsUserOwner(currentUserID, userID)
	if err != nil {
		return nil, err
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

	err = utility.IsUserOwner(currentUserID, userID)
	if err != nil {
		return err
	}

	deletedPost := postUseCase.postRepository.DeletePostByID(ctx, postID)

	return deletedPost
}
