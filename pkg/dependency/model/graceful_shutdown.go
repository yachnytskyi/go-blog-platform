package model

import (
	"context"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location  = "pkg.dependency.model."
	completed = "Graceful shutdown of the app"
)

func GracefulShutdown(ctx context.Context, logger Logger, closers ...Closer) {
	for _, closer := range closers {
		closer.Close(ctx)
	}

	logger.Info(domainError.NewInfoMessage(location+"GracefulShutdown", completed))
}
