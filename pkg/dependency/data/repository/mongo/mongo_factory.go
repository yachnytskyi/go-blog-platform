package mongo

import (
	"context"
	"fmt"

	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
)

type MongoDBFactory struct {
	MongoConfig config.MongoConfig
}

// NewRepository creates a new database instance.
func (f *MongoDBFactory) NewRepository(ctx context.Context) interface{} {
	mongoconn := options.Client().ApplyURI(f.MongoConfig.MongoURI)
	mongoClient, connectError := mongo.Connect(ctx, mongoconn)
	db := mongoClient.Database(f.MongoConfig.MongoDatabaseName)
	if validator.IsErrorNotNil(connectError) {
		panic(connectError)
	}
	connectError = mongoClient.Ping(ctx, readpref.Primary())
	if validator.IsErrorNotNil(connectError) {
		panic(connectError)
	}
	fmt.Println("Database successfully connected...")
	return db
}

// NewUserRepository creates a new UserRepository.
func (f *MongoDBFactory) NewUserRepository(db interface{}) user.UserRepository {
	mongoDB := db.(*mongo.Database)
	return userRepository.NewUserRepository(mongoDB)
}

// NewPostRepository creates a new PostRepository.
func (f *MongoDBFactory) NewPostRepository(db interface{}) post.PostRepository {
	mongoDB := db.(*mongo.Database)
	return postRepository.NewPostRepository(mongoDB)
}
