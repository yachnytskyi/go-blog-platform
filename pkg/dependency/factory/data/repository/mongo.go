package repository

import (
	"context"
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	postRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userRepository "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	location           = "pkg.dependency.data.repository.mongo."
	retryDelayInterval = 30 * time.Second
	maxRetryAttempts   = 5
)

type MongoDBRepository struct {
	Config      model.Config
	Logger      model.Logger
	MongoClient *mongo.Client
}

func NewMongoDBRepository(config model.Config, logger model.Logger) *MongoDBRepository {
	return &MongoDBRepository{
		Config: config,
		Logger: logger,
	}
}

func (mongoDBRepository *MongoDBRepository) NewRepository(ctx context.Context) any {
	config := mongoDBRepository.Config.GetConfig()
	var connectError error

	mongoConnection := options.Client().ApplyURI(config.MongoDB.URI)
	mongoClient := connectToMongo(ctx, location+"NewRepository", mongoDBRepository.Logger, mongoConnection)
	if validator.IsError(mongoClient.Error) {
		mongoDBRepository.Logger.Panic(domainError.NewInternalError(location+"NewRepository.mongoClient", constants.DatabaseConnectionFailure))
	}

	mongoDBRepository.MongoClient = mongoClient.Data
	connectError = pingMongo(ctx, location+"NewRepository", mongoDBRepository.Logger, mongoDBRepository.MongoClient)
	if validator.IsError(connectError) {
		mongoDBRepository.Logger.Panic(domainError.NewInternalError(location+"NewRepository.pingMongo", constants.DatabaseConnectionFailure))
	}

	mongoDBRepository.Logger.Info(domainError.NewInfoMessage(location+"NewRepository", constants.DatabaseConnectionSuccess))
	return mongoDBRepository.MongoClient.Database(config.MongoDB.Name)
}

func (mongoDBRepository MongoDBRepository) Close(ctx context.Context) {
	disconnectError := mongoDBRepository.MongoClient.Disconnect(ctx)
	if validator.IsError(disconnectError) {
		mongoDBRepository.Logger.Panic(domainError.NewInternalError(location+"Close.Disconnect", disconnectError.Error()))
	}

	mongoDBRepository.Logger.Info(domainError.NewInfoMessage(location+"Close", constants.DatabaseConnectionClosed))
}

func (mongoDBRepository MongoDBRepository) NewUserRepository(database any) user.UserRepository {
	mongoDB := database.(*mongo.Database)
	return userRepository.NewUserRepository(mongoDBRepository.Config, mongoDBRepository.Logger, mongoDB)
}

func (mongoDBRepository MongoDBRepository) NewPostRepository(database any) post.PostRepository {
	mongoDB := database.(*mongo.Database)
	return postRepository.NewPostRepository(mongoDBRepository.Logger, mongoDB)
}

// connectToMongo attempts to connect to MongoDB server with retries.
func connectToMongo(ctx context.Context, location string, logger model.Logger, mongoConnection *options.ClientOptions) common.Result[*mongo.Client] {
	var client *mongo.Client
	var connectError error
	var delay = time.Second

	for index := 0; index < maxRetryAttempts; index++ {
		client, connectError = mongo.Connect(ctx, mongoConnection)
		if connectError == nil {
			return common.NewResultOnSuccess[*mongo.Client](client)
		}

		logger.Error(domainError.NewInternalError(location+".connectToMongo.MongoClient.Connect", connectError.Error()))
		delay += retryDelayInterval
		time.Sleep(delay)
	}

	internalError := domainError.NewInternalError(location+".connectToMongo.MongoClient.Connect", connectError.Error())
	logger.Error(internalError)
	return common.NewResultOnFailure[*mongo.Client](internalError)
}

// pingMongo attempts to ping the MongoDB server with retries.
func pingMongo(ctx context.Context, location string, logger model.Logger, client *mongo.Client) error {
	var connectError error
	var delay = time.Second

	for index := 0; index < maxRetryAttempts; index++ {
		connectError = client.Ping(ctx, readpref.Primary())
		if connectError == nil {
			return nil
		}

		logger.Error(domainError.NewInternalError(location+".pingMongo.MongoClient.Ping", connectError.Error()))
		delay += retryDelayInterval
		time.Sleep(delay)
	}

	return connectError
}
