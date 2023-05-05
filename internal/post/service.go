package post

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Service interface {
	GetPostById(ctx context.Context, postID string) (*models.Post, error)
	GetAllPosts(ctx context.Context, page int, limit int) ([]*models.Post, error)
	CreatePost(ctx context.Context, user *models.PostCreate) (*models.Post, error)
	UpdatePostById(ctx context.Context, postID string, post *models.PostUpdate, currentUserID string) (*models.Post, error)
	DeletePostByID(ctx context.Context, postID string, currentUserID string) error
}
