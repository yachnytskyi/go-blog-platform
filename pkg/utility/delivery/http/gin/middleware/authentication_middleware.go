package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location      = "pkg.utility.delivery.http.gin.middleware."
	authorization = "Authorization"
	bearer        = "Bearer"
)

// AuthenticationMiddleware is a Gin middleware for handling user authentication using JWT tokens.
// Returns a Gin middleware handler function.
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Extract the access token from the request.
		accessToken, tokenError := extractAccessToken(ginContext)
		if validator.IsError(tokenError) {
			// Abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, tokenError, http.StatusUnauthorized)
			return
		}

		// Get the application configuration.
		accessTokenConfig := config.AppConfig.AccessToken
		// Validate the JWT token.
		userTokenPayload, validateAccessTokenError := domainUtility.ValidateJWTToken(accessToken, accessTokenConfig.PublicKey)
		if validator.IsError(validateAccessTokenError) {
			// Handle token validation error and respond with an unauthorized status and JSON error.
			httpAuthorizationError := httpError.NewHttpAuthorizationErrorView("", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, httpAuthorizationError, http.StatusUnauthorized)
			return
		}

		// Set user-related information in the Gin context for downstream handlers.
		ginContext.Set(constants.UserIDContext, userTokenPayload.UserID)
		ginContext.Set(constants.UserRoleContext, userTokenPayload.Role)

		// Continue to the next middleware or handler in the chain.
		ginContext.Next()
	}
}
