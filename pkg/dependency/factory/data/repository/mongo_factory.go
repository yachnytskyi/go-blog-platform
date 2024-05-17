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
	successfully_connected = "Database connection is established..."
	successfully_closed    = "Database connection has been successfully closed..."
	location               = "pkg.dependency.data.repository.mongo.NewRepository."
	unsupportedDatabase    = "Unsupported database type: %s"
)

// MongoDBFactory is a factory for creating MongoDB instances.
type MongoDBFactory struct {
	MongoDB     config.MongoDB // MongoDB configuration.
	MongoClient *mongo.Client  // MongoDB client instance.
}

// NewRepository creates and returns a new MongoDB repository instance.
func (mongoDBFactory *MongoDBFactory) NewRepository(ctx context.Context) any {
	var connectError error
	// Create a new MongoDB client with the provided URI.
	mongoConnection := options.Client().ApplyURI(mongoDBFactory.MongoDB.URI)
	mongoDBFactory.MongoClient, connectError = mongo.Connect(ctx, mongoConnection)
	if validator.IsError(connectError) {
		logging.Logger(domainError.NewInternalError(location+"MongoClient.Database", connectError.Error()))
		// Return nil to indicate the failure to establish a connection
		return nil
	}

	// Get the database instance from the MongoDB client.
	db := mongoDBFactory.MongoClient.Database(mongoDBFactory.MongoDB.Name)

	// Ping the MongoDB server to ensure a successful connection.
	connectError = mongoDBFactory.MongoClient.Ping(ctx, readpref.Primary())
	if validator.IsError(connectError) {
		logging.Logger(domainError.NewInternalError(location+"MongoClient.Ping", connectError.Error()))
		// Return nil to indicate the failure to establish a connection
		return nil
	}

	logging.Logger(successfully_connected)
	return db
}

// CloseRepository closes the MongoDB client and releases resources.
func (mongoDBFactory *MongoDBFactory) CloseRepository(ctx context.Context) {
	if validator.IsValueNotEmpty(mongoDBFactory.MongoClient) {
		mongoDBFactory.MongoClient.Disconnect(ctx)
		logging.Logger(successfully_closed)
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
