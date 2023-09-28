package domain

import (
	repository "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/data/repository"
	application "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/application"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

const (
	cleanup = "Cleaned the repository connection"
)

func GracefulShutdownDomain(databaseFactory repository.RepositoryFactory) {
	// Close any domain-specific resources or connections.
	// Add any additional domain-specific cleanup here.
	databaseFactory.CloseRepository()
	logging.Logger(cleanup)

	// Finally, perform application-level graceful shutdown.
	application.GracefulShutdown()
}
