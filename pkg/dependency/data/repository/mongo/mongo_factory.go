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

func (mongoDBFactory *MongoDBFactory) NewRepository(ctx context.Context) interface{} {
	mongoConnection := options.Client().ApplyURI(mongoDBFactory.MongoConfig.MongoURI)
	mongoClient, connectError := mongo.Connect(ctx, mongoConnection)
	defer mongoClient.Disconnect(ctx)
	db := mongoClient.Database(mongoDBFactory.MongoConfig.MongoDatabaseName)
	if validator.IsErrorNotNil(connectError) {
		logging.Logger(domainError.NewInternalError(location+"mongoClient.Database", connectError.Error()))
		return nil
	}
	connectError = mongoClient.Ping(ctx, readpref.Primary())
	if validator.IsErrorNotNil(connectError) {
		logging.Logger(domainError.NewInternalError(location+"mongoClient.Ping", connectError.Error()))
		return nil
	}
	fmt.Println("Database successfully connected...")
	return db
}

func (mongoDBFactory *MongoDBFactory) NewUserRepository(db interface{}) user.UserRepository {
	mongoDB := db.(*mongo.Database)
	return userRepository.NewUserRepository(mongoDB)
}

func (mongoDBFactory *MongoDBFactory) NewPostRepository(db interface{}) post.PostRepository {
	mongoDB := db.(*mongo.Database)
	return postRepository.NewPostRepository(mongoDB)
}
