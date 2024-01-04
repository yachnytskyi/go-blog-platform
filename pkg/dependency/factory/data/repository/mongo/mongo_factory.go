package mongo

import (
	"context"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
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
)

const (
	successfully_connected = "Database connection is established..."
	successfully_closed    = "Database connection has been successfully closed..."
	location               = "pkg.dependency.data.repository.mongo.NewRepository."
	unsupportedDatabase    = "Unsupported database type: %s"
)

type MongoDBFactory struct {
	MongoDB     config.MongoDB
	MongoClient *mongo.Client
}

func (mongoDBFactory *MongoDBFactory) NewRepository(ctx context.Context) any {
	var connectError error
	mongoConnection := options.Client().ApplyURI(mongoDBFactory.MongoDB.URI)
	mongoDBFactory.MongoClient, connectError = mongo.Connect(ctx, mongoConnection)
	db := mongoDBFactory.MongoClient.Database(mongoDBFactory.MongoDB.Name)
	if validator.IsError(connectError) {
		logging.Logger(domainError.NewInternalError(location+"mongoClient.Database", connectError.Error()))
		return nil
	}
	connectError = mongoDBFactory.MongoClient.Ping(ctx, readpref.Primary())
	if validator.IsError(connectError) {
		logging.Logger(domainError.NewInternalError(location+"mongoClient.Ping", connectError.Error()))
		return nil
	}
	logging.Logger(successfully_connected)
	return db
}

func (mongoDBFactory *MongoDBFactory) CloseRepository(ctx context.Context) {
	if validator.IsValueNotNil(mongoDBFactory.MongoClient) {
		mongoDBFactory.MongoClient.Disconnect(ctx)
		logging.Logger(successfully_closed)
	}
}

func (mongoDBFactory *MongoDBFactory) NewUserRepository(db any) user.UserRepository {
	mongoDB := db.(*mongo.Database)
	return userRepository.NewUserRepository(mongoDB)
}

func (mongoDBFactory *MongoDBFactory) NewPostRepository(db any) post.PostRepository {
	mongoDB := db.(*mongo.Database)
	return postRepository.NewPostRepository(mongoDB)
}
