package repository

import (
	"context"
	"time"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
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
	dbConnectionFailure = "Failed to establish database connection"
	retryDelayInterval  = 30 * time.Second
	maxRetryAttempts    = 5
)

// MongoDBRepository is a structure for creating MongoDB instances.
// It holds the client instance for connecting to MongoDB.
type MongoDBRepository struct {
	MongoClient *mongo.Client // MongoDB client instance.
}

func NewMongoDBRepository() *MongoDBRepository {
	return &MongoDBRepository{}
}

// NewRepository creates a new MongoDB repository instance based on the provided MongoDB configuration.
// It establishes a connection to the MongoDB server using retry logic and returns a handle to the MongoDB database.
// If connection or ping attempts fail after retries, it panics with a detailed error message.
func (mongoDBRepository *MongoDBRepository) NewRepository(ctx context.Context) any {
	// Load the MongoDB configuration.
	mongoDBConfig := config.GetMongoDBConfig()

	var connectError error

	// Configure MongoDB connection options using the provided URI.
	mongoConnection := options.Client().ApplyURI(mongoDBConfig.URI)

	// Try to connect to MongoDB with retries.
	mongoClient := connectToMongo(ctx, mongoConnection)

	// Assign the MongoDB client instance obtained from the successful connection result.
	mongoDBRepository.MongoClient = mongoClient.Data

	if validator.IsError(mongoClient.Error) {
		// Panic if all connection attempts fail.
		panic(dbConnectionFailure)
	}

	// Try to ping the MongoDB server with retries.
	connectError = pingMongo(ctx, mongoDBRepository.MongoClient)
	if validator.IsError(connectError) {
		// Panic if all ping attempts fail.
		panic(dbConnectionFailure)
	}

	// Log successful database connection.
	logging.Logger(dbConnectionSuccess)

	// Return the MongoDB database handle.
	return mongoDBRepository.MongoClient.Database(mongoDBConfig.Name)
}

// CloseRepository closes the MongoDB client and releases associated resources.
// It attempts to disconnect the MongoDB client with the provided context.
// Logs any errors encountered during disconnection and logs a success message upon successful closure.
//
// Parameters:
// - ctx: The context to control the timeout and cancellation for the disconnect operation.
func (mongoDBRepository *MongoDBRepository) CloseRepository(ctx context.Context) {
	// Attempt to disconnect the MongoDB client.
	disconnectError := mongoDBRepository.MongoClient.Disconnect(ctx)
	if validator.IsError(disconnectError) {
		// Log any errors that occur during disconnection.
		internalError := domainError.NewInternalError(location+"CloseRepository.Disconnect", disconnectError.Error())
		logging.Logger(internalError)
	}

	// Log a success message if the disconnection is successful.
	logging.Logger(dbConnectionClosed)
}

// NewUserRepository creates and returns a new UserRepository instance using the provided database.
//
// Parameters:
// - database: The MongoDB database instance.
//
// Returns:
// - user.UserRepository: The newly created UserRepository instance.
//
// This method type asserts the generic database to *mongo.Database and then
// creates a UserRepository with it.
func (mongoDBRepository *MongoDBRepository) NewUserRepository(database any) user.UserRepository {
	mongoDB := database.(*mongo.Database)
	return userRepository.NewUserRepository(mongoDB)
}

// NewPostRepository creates and returns a new PostRepository instance using the provided database.
//
// Parameters:
// - database: The MongoDB database instance.
//
// Returns:
// - post.PostRepository: The newly created PostRepository instance.
//
// This method type asserts the generic database to *mongo.Database and then
// creates a PostRepository with it.
func (mongoDBRepository *MongoDBRepository) NewPostRepository(database any) post.PostRepository {
	mongoDB := database.(*mongo.Database)
	return postRepository.NewPostRepository(mongoDB)
}

// connectToMongo attempts to connect to MongoDB server with retries.
//
// Parameters:
// - ctx: The context for controlling connection attempts.
// - mongoConnection: The options for connecting to MongoDB.
//
// Returns:
// - commonModel.Result[*mongo.Client]: The result containing the MongoDB client or an error.
//
// This function attempts to connect to MongoDB using exponential backoff.
// It retries the connection if it fails, up to a maximum number of attempts.
func connectToMongo(ctx context.Context, mongoConnection *options.ClientOptions) commonModel.Result[*mongo.Client] {
	var client *mongo.Client
	var connectError error
	var delay = time.Second

	// Attempt to connect to MongoDB with exponential backoff.
	for index := 0; index < maxRetryAttempts; index++ {
		client, connectError = mongo.Connect(ctx, mongoConnection)
		if connectError == nil {
			// Return client if connection is successful.
			return commonModel.NewResultOnSuccess[*mongo.Client](client)
		}

		// Log the connection error with detailed message and retry after delay.
		logging.Logger(domainError.NewInternalError(location+"connectToMongo.MongoClient.Connect", connectError.Error()))
		delay += retryDelayInterval
		time.Sleep(delay)
	}

	// Log the connection error with detailed message.
	logging.Logger(domainError.NewInternalError(location+"connectToMongo.MongoClient.Connect", connectError.Error()))

	// Return error if all retry attempts fail.
	return commonModel.NewResultOnFailure[*mongo.Client](connectError)
}

// pingMongo attempts to ping the MongoDB server with retries.
//
// Parameters:
// - ctx: The context for controlling ping attempts.
// - client: The MongoDB client instance.
//
// Returns:
// - error: The error if pinging fails after all attempts, otherwise nil.
//
// This function attempts to ping the MongoDB server using exponential backoff.
// It retries the ping if it fails, up to a maximum number of attempts.
func pingMongo(ctx context.Context, client *mongo.Client) error {
	var connectError error
	var delay = time.Second

	// Attempt to ping MongoDB server with exponential backoff.
	for index := 0; index < maxRetryAttempts; index++ {
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
