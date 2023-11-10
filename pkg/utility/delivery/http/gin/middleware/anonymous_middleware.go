package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// AnonymousMiddleware is a Gin middleware to check if the user is anonymous based on the presence of an access token.
func AnonymousMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the access token from the request.
		accessToken, tokenError := extractAccessToken(ginContext)
		if tokenError != nil {
			// Abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, tokenError, http.StatusUnauthorized)
			return
		}

		// Check for a deadline error using the handleDeadlineExceeded function.
		// If a deadline error occurred, respond with a timeout status.
		deadlineError := handleDeadlineExceeded(ctx)
		if validator.IsValueNotNil(deadlineError) {
			// Use the abortWithStatusJSON function to handle the deadline error by sending
			// a JSON response with an appropriate HTTP status code.
			abortWithStatusJSON(ginContext, deadlineError, http.StatusUnauthorized)
		}

		// Check if the access token is not empty, indicating that the user is already authenticated.
		if validator.IsStringNotEmpty(accessToken) {
			// Create a custom error message indicating that the user is already authenticated.
			authorizationError := httpError.NewHttpAuthorizationErrorView(location+"AnonymousMiddleware.accessToken", constants.AlreadyRegisteredNotification)
			logging.Logger(authorizationError)
			jsonResponse := httpModel.NewJsonResponseOnFailure(authorizationError)
			ginContext.AbortWithStatusJSON(http.StatusForbidden, jsonResponse)
			return
		}

		// Continue to the next middleware or handler in the chain.
		ginContext.Next()
	}
}
