package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/interfaces"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
)

// AnonymousMiddleware is a Gin middleware to check if the user is anonymous based on the presence of an access token.
func AnonymousMiddleware(logger interfaces.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		anonymousAccessToken := isUserAnonymous(ginContext)

		// If the access token is present, indicating that the user is already authenticated.
		if anonymousAccessToken {
			httpAuthorizationError := delivery.NewHTTPAuthorizationError(location+"AnonymousMiddleware.anonymousAccessToken", constants.AlreadyLoggedInNotification)
			abortWithStatusJSON(ginContext, logger, httpAuthorizationError, http.StatusForbidden)
			return
		}

		ginContext.Next()
	}
}

func isUserAnonymous(ginContext *gin.Context) bool {
	authorizationHeader := ginContext.Request.Header.Get(constants.Authorization)

	// Attempt to retrieve the access token from the cookie if no valid Bearer token is found in the Authorization header.
	_, cookieError := ginContext.Cookie(constants.AccessTokenValue)
	if cookieError == nil {
		return true
	}

	// Verify that the Authorization header contains a Bearer token.
	if strings.HasPrefix(authorizationHeader, constants.Bearer) && len(authorizationHeader) > len(constants.Bearer) {
		return true
	}

	return false
}
