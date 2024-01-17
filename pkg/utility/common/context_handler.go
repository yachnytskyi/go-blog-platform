package common

import (
	"context"

	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
)

// HandleWithContextError is a utility function that encapsulates
// the common logic for context error handling.
func HandleWithContextError(location string, ctx context.Context) error {
	// Check for context cancellation after serving the request.
	select {
	case <-ctx.Done():
		internalError := domainError.NewInternalError(location+".HandleWithContextError", ctx.Err().Error())

		// Log the context error.
		logging.Logger(internalError)

		// Handle the error using the provided errorHandler.
		return domainError.HandleError(internalError)

		// Handle the error using your custom logic.
		// You can add more logic based on your needs.
	default:
		// No context cancellation.
		// Continue with normal execution.
		return nil
	}
}
