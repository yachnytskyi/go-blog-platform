package post

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Repository interface {
	GetPostById(ctx context.Context, postID string) (*models.PostDB, error)
	GetAllPosts(ctx context.Context, page int, limit int) ([]*models.PostDB, error)
	CreatePost(ctx context.Context, user *models.PostCreate) (*models.PostDB, error)
	UpdatePostById(ctx context.Context, postID string, post *models.PostUpdate) (*models.PostDB, error)
	DeletePostByID(ctx context.Context, postID string) error
}
