package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
)

// AnonymousMiddleware is a Gin middleware to check if the user is anonymous based on the presence of an access token.
func AnonymousMiddleware(logger model.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		anonymousAccessToken := isUserAnonymous(ginContext)

		// If the access token is present, indicating that the user is already authenticated.
		if anonymousAccessToken {
			httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+"AnonymousMiddleware.anonymousAccessToken", constants.AlreadyLoggedInNotification)
			abortWithStatusJSON(ginContext, logger, httpAuthorizationError, constants.StatusForbidden)
			return
		}

		ginContext.Next()
	}
}

// isUserAnonymous checks if the user is anonymous based on the presence of an access token.
func isUserAnonymous(ginContext *gin.Context) bool {
	authorizationHeader := ginContext.Request.Header.Get(constants.Authorization)
	fields := strings.Fields(authorizationHeader)

	// Check if the Authorization header contains a Bearer token.
	if len(fields) > 0 && fields[0] == constants.Bearer {
		return true
	}

	// If no Bearer token in the Authorization header, try to get the token from the cookie.
	_, cookieError := ginContext.Cookie(constants.AccessTokenValue)
	if cookieError == nil {
		return true
	}

	return false
}
