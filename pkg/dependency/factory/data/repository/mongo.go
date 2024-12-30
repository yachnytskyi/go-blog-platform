package repository

import (
	"context"
	"fmt"
	"time"

	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/data/repository/mongo"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/data/repository/mongo"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/interfaces"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	location                           = "pkg.dependency.factory.data.repository."
	connectingToMongoDBNotification    = "Attempting to connect to the MongoDB database..."
	databasePingingMongoDBNotification = "Attempting to ping the MongoDB database..."

	retryDelayInterval        = 1 * time.Second
	maxRetryAttempts          = 2
	healthCheckTickerInterval = 60 * time.Second
)

type MongoDBRepository struct {
	Config      *config.ApplicationConfig
	Logger      interfaces.Logger
	MongoClient *mongo.Client
}

func NewMongoDBRepository(config *config.ApplicationConfig, logger interfaces.Logger) *MongoDBRepository {
	return &MongoDBRepository{
		Config: config,
		Logger: logger,
	}
}

func (mongoDBRepository *MongoDBRepository) CreateRepository(ctx context.Context) any {
	mongoDBRepository.connectToMongoDB(ctx)
	return mongoDBRepository.MongoClient.Database(mongoDBRepository.Config.MongoDB.Name)
}

func (mongoDBRepository MongoDBRepository) NewRepository(createRepository any, repository any) any {
	mongoDB := createRepository.(*mongo.Database)
	switch repository.(type) {
	case *interfaces.UserRepository:
		return user.NewUserRepository(mongoDBRepository.Config, mongoDBRepository.Logger, mongoDB)
	case *interfaces.PostRepository:
		return post.NewPostRepository(mongoDBRepository.Logger, mongoDB)
	default:
		mongoDBRepository.Logger.Panic(domain.NewInternalError(location+"mongo.NewRepository.default", fmt.Sprintf(constants.UnsupportedRepository, repository)))
		return nil
	}
}

func (mongoDBRepository *MongoDBRepository) HealthCheck(delivery interfaces.Delivery) {
	go func() {
		ticker := time.NewTicker(healthCheckTickerInterval)
		defer ticker.Stop()

		for range ticker.C {
			databasePing := mongoDBRepository.DatabasePing()
			if !databasePing {
				mongoDBRepository.handleReconnection(delivery)
			}
		}
	}()
}

func (mongoDBRepository *MongoDBRepository) DatabasePing() bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	pingError := mongoDBRepository.MongoClient.Ping(ctx, readpref.Primary())
	return pingError == nil
}

func (mongoDBRepository MongoDBRepository) Close(ctx context.Context) {
	disconnectError := mongoDBRepository.MongoClient.Disconnect(ctx)
	if validator.IsError(disconnectError) {
		mongoDBRepository.Logger.Panic(domain.NewInternalError(location+"mongo.Close.Disconnect", disconnectError.Error()))
	}

	mongoDBRepository.Logger.Info(domain.NewInfoMessage(location+"mongo.Close", constants.DatabaseConnectionClosed))
}

func (mongoDBRepository *MongoDBRepository) connectToMongoDB(ctx context.Context) {
	var client *mongo.Client
	var connectError error
	var delay = time.Second

	for index := 0; index < maxRetryAttempts; index++ {
		mongoDBRepository.Logger.Warn(domain.NewInfoMessage(location+"mongo.connectToMongoDB", connectingToMongoDBNotification))
		client, connectError = mongo.Connect(ctx, options.Client().ApplyURI(mongoDBRepository.Config.MongoDB.URI))
		if connectError == nil {
			pingError := client.Ping(ctx, readpref.Primary())
			if pingError == nil {
				mongoDBRepository.MongoClient = client
				mongoDBRepository.Logger.Info(domain.NewInfoMessage(location+"mongo.connectToMongoDB", constants.DatabaseConnectionSuccess))
				return
			}
		}
		time.Sleep(delay)
		delay += retryDelayInterval
	}

	mongoDBRepository.Logger.Panic(domain.NewInternalError(location+"mongo.connectToMongoDB", connectError.Error()))
}

func (mongoDBRepository *MongoDBRepository) handleReconnection(delivery interfaces.Delivery) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	delay := time.Second
	for index := 0; index < maxRetryAttempts; index++ {
		mongoDBRepository.Logger.Warn(domain.NewInfoMessage(location+"mongo.handleReconnection", connectingToMongoDBNotification))
		client, connectError := mongo.Connect(ctx, options.Client().ApplyURI(mongoDBRepository.Config.MongoDB.URI))
		if connectError == nil {
			databasePing := mongoDBRepository.DatabasePing()
			if databasePing {
				mongoDBRepository.MongoClient = client
				mongoDBRepository.Logger.Info(domain.NewInfoMessage(location+"mongo.handleReconnection", constants.DatabaseConnectionSuccess))
				return
			}
		}
		time.Sleep(delay)
		delay += retryDelayInterval
	}

	model.GracefulShutdown(ctx, mongoDBRepository.Logger, mongoDBRepository, delivery)
	mongoDBRepository.Logger.Panic(domain.NewInternalError(location+"mongo.handleReconnection", constants.DatabaseConnectionFailure))
}
