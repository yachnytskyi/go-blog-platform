package model

import (
	"context"

	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
)

const (
	location  = "pkg.dependency.model."
	completed = "Graceful shutdown of the app"
)

func GracefulShutdown(ctx context.Context, logger interfaces.Logger, close ...interfaces.Close) {
	for _, close := range close {
		close.Close(ctx)
	}

	logger.Info(domainError.NewInfoMessage(location+"GracefulShutdown", completed))
}
