package repository

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	mongoDBFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository/mongo"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location            = "pkg/dependency/data/repository/InjectRepository"
	unsupportedDatabase = "unsupported database type: %s"
)

// Define a DatabaseFactory interface to create different database instances.
type DatabaseFactory interface {
	NewRepository(ctx context.Context) (interface{}, error)
	NewUserRepository(db interface{}) user.UserRepository
	NewPostRepository(db interface{}) post.PostRepository
}

func InjectRepository(loadConfig config.Config) (DatabaseFactory, error) {
	switch loadConfig.Database {
	case loadConfig.MongoConfig:
		return &mongoDBFactory.MongoDBFactory{MongoConfig: loadConfig.MongoConfig}, nil
	// Add other database cases here as needed.
	default:
		internalError := domainError.NewInternalError(location+".loadConfig.MongoConfig:", fmt.Sprintf(unsupportedDatabase, loadConfig.Database))
		logging.Logger(internalError)
		return nil, domainError.NewInternalError(location, fmt.Sprintf(unsupportedDatabase, loadConfig.Database))
	}
}
