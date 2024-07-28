package repository

import (
	"context"
	"fmt"
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	location               = "pkg.dependency.factory.data.repository.mongo."
	retryDelayInterval     = 30 * time.Second
	maxRetryAttempts       = 5
	repositoryDoesNotExist = "This repository type does not exist: %s"
)

type MongoDBRepository struct {
	Config      interfaces.Config
	Logger      interfaces.Logger
	MongoClient *mongo.Client
}

func NewMongoDBRepository(config interfaces.Config, logger interfaces.Logger) *MongoDBRepository {
	return &MongoDBRepository{
		Config: config,
		Logger: logger,
	}
}

func (mongoDBRepository *MongoDBRepository) CreateRepository(ctx context.Context) any {
	config := mongoDBRepository.Config.GetConfig()
	var connectError error

	mongoConnection := options.Client().ApplyURI(config.MongoDB.URI)
	mongoClient := connectToMongo(ctx, location+"CreateRepository", mongoDBRepository.Logger, mongoConnection)
	if validator.IsError(mongoClient.Error) {
		mongoDBRepository.Logger.Panic(domainError.NewInternalError(location+"CreateRepository.mongoClient", constants.DatabaseConnectionFailure))
	}

	mongoDBRepository.MongoClient = mongoClient.Data
	connectError = pingMongo(ctx, location+"CreateRepository", mongoDBRepository.Logger, mongoDBRepository.MongoClient)
	if validator.IsError(connectError) {
		mongoDBRepository.Logger.Panic(domainError.NewInternalError(location+"CreateRepository.pingMongo", constants.DatabaseConnectionFailure))
	}

	mongoDBRepository.Logger.Info(domainError.NewInfoMessage(location+"CreateRepository", constants.DatabaseConnectionSuccess))
	return mongoDBRepository.MongoClient.Database(config.MongoDB.Name)
}

func (mongoDBRepository MongoDBRepository) NewRepository(createRepository any, repository any) any {
	mongoDB := createRepository.(*mongo.Database)
	switch repository.(type) {
	case *interfaces.UserRepository:
		return user.NewUserRepository(mongoDBRepository.Config, mongoDBRepository.Logger, mongoDB)
	case *interfaces.PostRepository:
		return post.NewPostRepository(mongoDBRepository.Logger, mongoDB)
	default:
		mongoDBRepository.Logger.Panic(domainError.NewInternalError(location+"NewRepository.default", fmt.Sprintf(repositoryDoesNotExist, repository)))
		return nil
	}
}

func (mongoDBRepository MongoDBRepository) Close(ctx context.Context) {
	disconnectError := mongoDBRepository.MongoClient.Disconnect(ctx)
	if validator.IsError(disconnectError) {
		mongoDBRepository.Logger.Panic(domainError.NewInternalError(location+"Close.Disconnect", disconnectError.Error()))
	}

	mongoDBRepository.Logger.Info(domainError.NewInfoMessage(location+"Close", constants.DatabaseConnectionClosed))
}

// connectToMongo attempts to connect to MongoDB server with retries.
func connectToMongo(ctx context.Context, location string, logger interfaces.Logger, mongoConnection *options.ClientOptions) common.Result[*mongo.Client] {
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
func pingMongo(ctx context.Context, location string, logger interfaces.Logger, client *mongo.Client) error {
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
