package mongo

import (
	"context"
	"fmt"

	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
)

const (
	location            = "pkg.dependency.data.repository.mongo.NewRepoitory"
	unsupportedDatabase = "unsupported database type: %s"
)

type MongoDBFactory struct {
	MongoConfig config.MongoConfig
}

// NewRepository creates a new database instance.
func (f *MongoDBFactory) NewRepository(ctx context.Context) (interface{}, error) {
	mongoconn := options.Client().ApplyURI(f.MongoConfig.MongoURI)
	mongoClient, connectError := mongo.Connect(ctx, mongoconn)
	db := mongoClient.Database(f.MongoConfig.MongoDatabaseName)
	if validator.IsErrorNotNil(connectError) {
		internalError := domainError.NewInternalError(location+"mongoClient.Database", connectError.Error())
		logging.Logger(internalError)
		return nil, internalError
	}
	connectError = mongoClient.Ping(ctx, readpref.Primary())
	if validator.IsErrorNotNil(connectError) {
		internalError := domainError.NewInternalError(location+"mongoClient.Ping", connectError.Error())
		logging.Logger(internalError)
		return nil, internalError
	}
	fmt.Println("Database successfully connected...")
	return db, nil
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
