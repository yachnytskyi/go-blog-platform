package repository

import (
	"context"
	"time"

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
	retryDelayInterval  = 30 * time.Second
	maxRetryAttempts    = 5
)

// MongoDBFactory is a factory for creating MongoDB instances.
type MongoDBFactory struct {
	MongoDB     config.MongoDB // MongoDB configuration.
	MongoClient *mongo.Client  // MongoDB client instance.
}

// NewRepository creates and returns a new MongoDB repository instance.
func (mongoDBFactory *MongoDBFactory) NewRepository(ctx context.Context) any {
	var connectError error

	// Configure MongoDB connection options using the provided URI.
	mongoConnection := options.Client().ApplyURI(mongoDBFactory.MongoDB.URI)

	// Try to connect to MongoDB with retries.
	mongoDBFactory.MongoClient, connectError = connectToMongo(ctx, mongoConnection)
	if validator.IsError(connectError) {
		// Panic if all connection attempts fail.
		panic(dbConnectionFailure)
	}

	// Try to ping the MongoDB server with retries.
	connectError = pingMongo(ctx, mongoDBFactory.MongoClient)
	if validator.IsError(connectError) {
		// Panic if all connection attempts fail.
		panic(dbConnectionFailure)
	}

	// Log successful database connection.
	logging.Logger(dbConnectionSuccess)
	return mongoDBFactory.MongoClient.Database(mongoDBFactory.MongoDB.Name)
}

// CloseRepository closes the MongoDB client and releases resources.
// Parameters:
// - ctx: The context to control the timeout and cancellation for the disconnect operation.
func (mongoDBFactory *MongoDBFactory) CloseRepository(ctx context.Context) {
	// Attempt to disconnect the MongoDB client.
	disconnectError := mongoDBFactory.MongoClient.Disconnect(ctx)
	if validator.IsError(disconnectError) {
		// Log any errors that occur during disconnection.
		internalError := domainError.NewInternalError(location+"CloseRepository.Disconnect", disconnectError.Error())
		logging.Logger(internalError)
	}

	// Log a success message if the disconnection is successful.
	logging.Logger(dbConnectionClosed)
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

// connectToMongo attempts to connect to MongoDB server with retries.
func connectToMongo(ctx context.Context, mongoConnection *options.ClientOptions) (*mongo.Client, error) {
	var client *mongo.Client
	var connectError error
	var delay = time.Second

	// Attempt to connect to MongoDB with exponential backoff.
	for i := 0; i < maxRetryAttempts; i++ {
		client, connectError = mongo.Connect(ctx, mongoConnection)
		if connectError == nil {
			// Return client if connection is successful.
			return client, nil
		}

		// Log the connection error with detailed message and retry after delay.
		logging.Logger(domainError.NewInternalError(location+"connectToMongo.MongoClient.Connect", connectError.Error()))
		delay += retryDelayInterval
		time.Sleep(delay)
	}

	// Return error if all retry attempts fail.
	return nil, connectError
}

// pingMongo attempts to ping the MongoDB server with retries.
func pingMongo(ctx context.Context, client *mongo.Client) error {
	var connectError error
	var delay = time.Second

	// Attempt to ping MongoDB server with exponential backoff.
	for i := 0; i < maxRetryAttempts; i++ {
		connectError = client.Ping(ctx, readpref.Primary())
		if connectError == nil {
			// Return nil if ping is successful.
			return nil
		}

		// Log the ping error with detailed message and retry after delay.
		logging.Logger(domainError.NewInternalError(location+"pingMongo.MongoClient.Ping", connectError.Error()))
		delay += retryDelayInterval
		time.Sleep(delay)
	}

	// Return error if all retry attempts fail.
	return connectError
}
