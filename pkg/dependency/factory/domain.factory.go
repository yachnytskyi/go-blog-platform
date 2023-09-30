package factory

import (
	"github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
)

// Define a DatabaseFactory interface to create different database instances.
type DomainFactory interface {
	NewUserRepository(db interface{}) user.UserUseCase
	NewPostRepository(db interface{}) post.PostUseCase
}
