package factory

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
)

// Define a DatabaseFactory interface to create different database instances.
type RepositoryFactory interface {
	NewRepository(ctx context.Context) interface{}
	CloseRepository()
	NewUserRepository(db interface{}) user.UserRepository
	NewPostRepository(db interface{}) post.PostRepository
}
