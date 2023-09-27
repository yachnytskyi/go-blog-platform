package repository

import (
	"context"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	mongoDBFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository/mongo"
)

// Define a DatabaseFactory interface to create different database instances.
type DatabaseFactory interface {
	NewRepository(ctx context.Context) interface{}
	NewUserRepository(db interface{}) user.UserRepository
	NewPostRepository(db interface{}) post.PostRepository
}

func InjectRepository(loadConfig config.Config) DatabaseFactory {
	if loadConfig.Database == loadConfig.MongoConfig {
		return &mongoDBFactory.MongoDBFactory{MongoConfig: loadConfig.MongoConfig}
	}
	// Add other database factories here as needed.
	return &mongoDBFactory.MongoDBFactory{MongoConfig: loadConfig.MongoConfig}
}
