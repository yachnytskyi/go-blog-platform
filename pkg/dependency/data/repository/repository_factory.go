package repository

import (
	"context"
	"fmt"

	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	mongoDBFactory "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository/mongo"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/application"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	location            = "pkg/dependency/data/repository/InjectRepository"
	unsupportedDatabase = "unsupported database type: %s"
)

// Define a DatabaseFactory interface to create different database instances.
type DatabaseFactory interface {
	NewRepository(ctx context.Context) interface{}
	NewUserRepository(db interface{}) user.UserRepository
	NewPostRepository(db interface{}) post.PostRepository
}

func InjectRepository(loadConfig config.Config) DatabaseFactory {
	switch loadConfig.Database {
	case config.MongoDB:
		return &mongoDBFactory.MongoDBFactory{MongoConfig: loadConfig.MongoConfig}
	// Add other database cases here as needed.
	default:
		logging.Logger(domainError.NewInternalError(location+".loadConfig.DatabaseType:", fmt.Sprintf(unsupportedDatabase, loadConfig.Database)))
		application.GracefulShutdown()
		return nil
	}
}
