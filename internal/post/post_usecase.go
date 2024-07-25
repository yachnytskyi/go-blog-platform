package post

import (
	"context"

	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
)

type PostUseCase interface {
	GetAllPosts(ctx context.Context, page int, limit int) (*post.Posts, error)
	GetPostById(ctx context.Context, postID string) (*post.Post, error)
	CreatePost(ctx context.Context, user *post.PostCreate) (*post.Post, error)
	UpdatePostById(ctx context.Context, postID string, post *post.PostUpdate, currentUserID string) (*post.Post, error)
	DeletePostByID(ctx context.Context, postID string, currentUserID string) error
}
