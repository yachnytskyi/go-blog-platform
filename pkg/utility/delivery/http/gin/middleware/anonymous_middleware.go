package middleware

import (
	"context"
	"strings"

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

		// Check if a user is anonymous
		anonymousAccessToken := isUserAnonymous(ginContext)

		// Check for a deadline error using the handleDeadlineExceeded function.
		// If a deadline error occurred, respond with a timeout status.
		deadlineError := handleDeadlineExceeded(ctx)
		if validator.IsValueNotNil(deadlineError) {
			// Use the abortWithStatusJSON function to handle the deadline error by sending
			// a JSON response with an appropriate HTTP status code.
			abortWithStatusJSON(ginContext, deadlineError, constants.StatusUnauthorized)
		}

		// Check if the access token is not empty, indicating that the user is already authenticated.
		if validator.IsStringNotEmpty(anonymousAccessToken) {
			// Create a custom error message indicating that the user is already authenticated.
			authorizationError := httpError.NewHttpAuthorizationErrorView(location+"AnonymousMiddleware.anonymousAccessToken", constants.AlreadyRegisteredNotification)
			logging.Logger(authorizationError)
			jsonResponse := httpModel.NewJSONResponseOnFailure(authorizationError)
			ginContext.AbortWithStatusJSON(constants.StatusForbidden, jsonResponse)
			return
		}

		// Continue to the next middleware or handler in the chain.
		ginContext.Next()
	}
}

// isUserAnonymous checks if the user is anonymous and returns the access token if present.
func isUserAnonymous(ginContext *gin.Context) string {
	var anonymousAccessToken string

	// Attempt to retrieve the access token from the Authorization header.
	authorizationHeader := ginContext.Request.Header.Get(authorization)
	fields := strings.Fields(authorizationHeader)

	// Check if the Authorization header contains a Bearer token.
	if validator.IsSliceNotEmpty(fields) && fields[firstElement] == bearer {
		anonymousAccessToken = fields[nextElement]
	} else {
		// If no Bearer token in the Authorization header, try to get the token from the cookie.
		cookie, cookieError := ginContext.Cookie(constants.AccessTokenValue)
		if cookieError == nil {
			anonymousAccessToken = cookie
		}
	}
	return anonymousAccessToken
}
