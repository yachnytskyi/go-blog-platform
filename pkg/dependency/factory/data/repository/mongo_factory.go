package repository

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
	dbConnectionSuccess = "Database connection is established..."
	dbConnectionClosed  = "Database connection has been successfully closed..."
	location            = "pkg.dependency.data.repository.mongo.NewRepository."
	unsupportedDatabase = "Unsupported database type: %s"
	dbConnectionFailure = "Failed to establish database connection"
)

// MongoDBFactory is a factory for creating MongoDB instances.
type MongoDBFactory struct {
	MongoDB     config.MongoDB // MongoDB configuration.
	MongoClient *mongo.Client  // MongoDB client instance.
}

// NewRepository creates and returns a new MongoDB repository instance.
func (mongoDBFactory *MongoDBFactory) NewRepository(ctx context.Context) any {
	var connectError error

	// Attempt to connect to the MongoDB server using the provided URI.
	mongoConnection := options.Client().ApplyURI(mongoDBFactory.MongoDB.URI)
	mongoDBFactory.MongoClient, connectError = mongo.Connect(ctx, mongoConnection)
	if validator.IsError(connectError) {
		// Log the connection error with a detailed message indicating the failure location.
		logging.Logger(domainError.NewInternalError(location+"MongoClient.Connect", connectError.Error()))
		panic(dbConnectionFailure)
	}

	// Ping the MongoDB server to ensure a successful connection.
	connectError = mongoDBFactory.MongoClient.Ping(ctx, readpref.Primary())
	if validator.IsError(connectError) {
		// Log the ping error with a detailed message indicating the failure location.
		logging.Logger(domainError.NewInternalError(location+"MongoClient.Ping", connectError.Error()))
		// Panic to stop execution if the MongoDB server cannot be reached after the initial connection.
		panic(dbConnectionFailure)
	}

	logging.Logger(dbConnectionSuccess)
	return mongoDBFactory.MongoClient.Database(mongoDBFactory.MongoDB.Name)
}

// CloseRepository closes the MongoDB client and releases resources.
func (mongoDBFactory *MongoDBFactory) CloseRepository(ctx context.Context) {
	if validator.IsValueNotEmpty(mongoDBFactory.MongoClient) {
		mongoDBFactory.MongoClient.Disconnect(ctx)
		logging.Logger(dbConnectionClosed)
	}
}

// NewUserRepository creates and returns a new UserRepository instance using the provided database.
func (mongoDBFactory *MongoDBFactory) NewUserRepository(db any) user.UserRepository {
	mongoDB := db.(*mongo.Database)
	return userRepository.NewUserRepository(mongoDB)
}

// NewPostRepository creates and returns a new PostRepository instance using the provided database.
func (mongoDBFactory *MongoDBFactory) NewPostRepository(db any) post.PostRepository {
	mongoDB := db.(*mongo.Database)
	return postRepository.NewPostRepository(mongoDB)
}
