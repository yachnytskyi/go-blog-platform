package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// AnonymousMiddleware is a Gin middleware to check if the user is anonymous based on the presence of an access token.
func AnonymousMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Check if a user is anonymous
		anonymousAccessToken := isUserAnonymous(ginContext)

		// Check if the access token is not empty, indicating that the user is already authenticated.
		if validator.IsStringNotEmpty(anonymousAccessToken) {
			// Create a custom error message indicating that the user is already authenticated.
			authorizationError := httpError.NewHttpAuthorizationErrorView(location+"AnonymousMiddleware.anonymousAccessToken", constants.AlreadyLoggedInNotification)
			abortWithStatusJSON(ginContext, authorizationError, constants.StatusForbidden)
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
