package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
)

// TimeoutMiddleware sets a timeout for each request.
// Returns a Gin middleware handler function.
func TimeoutMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Create a context with timeout.
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Create a new request with the modified context.
		ginContext.Request = ginContext.Request.WithContext(ctx)

		// Use a goroutine to call Next and handle the response.
		ch := make(chan struct{})
		go func() {
			defer close(ch)
			ginContext.Next()
		}()

		// Use a select statement to wait for either the request to complete or the context to timeout.
		select {
		case <-ch:
			// Request completed successfully.
		case <-ctx.Done():
			// Context timed out.
			httpInternalErrorView := httpError.NewHttpInternalErrorView(location+"TimeOutMiddleware", ctx.Err().Error())
			// Abort the request with an HTTP status and respond with a JSON error.
			abortWithStatusJSON(ginContext, httpInternalErrorView, constants.StatusBadGateway)
		}
	}
}
