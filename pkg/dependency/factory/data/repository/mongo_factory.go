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

type MongoDBRepository struct {
	MongoClient *mongo.Client
}

func NewMongoDBRepository() *MongoDBRepository {
	return &MongoDBRepository{}
}

func (mongoDBRepository *MongoDBRepository) NewRepository(ctx context.Context) any {
	mongoDBConfig := config.GetMongoDBConfig()
	var connectError error

	mongoConnection := options.Client().ApplyURI(mongoDBConfig.URI)
	mongoClient := connectToMongo(ctx, mongoConnection)
	mongoDBRepository.MongoClient = mongoClient.Data
	if validator.IsError(mongoClient.Error) {
		panic(dbConnectionFailure)
	}

	connectError = pingMongo(ctx, mongoDBRepository.MongoClient)
	if validator.IsError(connectError) {
		panic(dbConnectionFailure)
	}

	logging.Logger(dbConnectionSuccess)
	return mongoDBRepository.MongoClient.Database(mongoDBConfig.Name)
}

func (mongoDBRepository *MongoDBRepository) CloseRepository(ctx context.Context) {
	disconnectError := mongoDBRepository.MongoClient.Disconnect(ctx)
	if validator.IsError(disconnectError) {
		internalError := domainError.NewInternalError(location+"CloseRepository.Disconnect", disconnectError.Error())
		logging.Logger(internalError)
	}

	logging.Logger(dbConnectionClosed)
}

func (mongoDBRepository *MongoDBRepository) NewUserRepository(database any) user.UserRepository {
	mongoDB := database.(*mongo.Database)
	return userRepository.NewUserRepository(mongoDB)
}

func (mongoDBRepository *MongoDBRepository) NewPostRepository(database any) post.PostRepository {
	mongoDB := database.(*mongo.Database)
	return postRepository.NewPostRepository(mongoDB)
}

// connectToMongo attempts to connect to MongoDB server with retries.
func connectToMongo(ctx context.Context, mongoConnection *options.ClientOptions) commonModel.Result[*mongo.Client] {
	var client *mongo.Client
	var connectError error
	var delay = time.Second

	for index := 0; index < maxRetryAttempts; index++ {
		client, connectError = mongo.Connect(ctx, mongoConnection)
		if connectError == nil {
			return commonModel.NewResultOnSuccess[*mongo.Client](client)
		}

		logging.Logger(domainError.NewInternalError(location+"connectToMongo.MongoClient.Connect", connectError.Error()))
		delay += retryDelayInterval
		time.Sleep(delay)
	}

	internalError := domainError.NewInternalError(location+"connectToMongo.MongoClient.Connect", connectError.Error())
	logging.Logger(internalError)
	return commonModel.NewResultOnFailure[*mongo.Client](internalError)
}

// pingMongo attempts to ping the MongoDB server with retries.
func pingMongo(ctx context.Context, client *mongo.Client) error {
	var connectError error
	var delay = time.Second

	for index := 0; index < maxRetryAttempts; index++ {
		connectError = client.Ping(ctx, readpref.Primary())
		if connectError == nil {
			return nil
		}

		logging.Logger(domainError.NewInternalError(location+"pingMongo.MongoClient.Ping", connectError.Error()))
		delay += retryDelayInterval
		time.Sleep(delay)
	}

	return connectError
}
