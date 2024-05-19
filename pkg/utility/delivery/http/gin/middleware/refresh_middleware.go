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

// RefreshTokenAuthenticationMiddleware is a Gin middleware for handling user authentication using refresh tokens.
// Returns a Gin middleware handler function.
func RefreshTokenAuthenticationMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Extract the refresh token from the request.
		refreshToken, tokenError := extractRefreshToken(ginContext)
		if validator.IsError(tokenError) {
			// Abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, tokenError, http.StatusUnauthorized)
			return
		}

		// Get the refresh token configuration.
		applicationConfig := config.AppConfig.RefreshToken

		// Validate the JWT token.
		userTokenPayload, validateRefreshTokenError := domainUtility.ValidateJWTToken(refreshToken, applicationConfig.PublicKey)
		if validator.IsError(validateRefreshTokenError) {
			// Handle token validation error and respond with an unauthorized status and JSON error.
			httpAuthorizationError := httpError.NewHttpAuthorizationErrorView("", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, httpAuthorizationError, http.StatusUnauthorized)
			return
		}

		// Set user-related information in the Gin context for downstream handlers.
		ginContext.Set(constants.UserIDContext, userTokenPayload.UserID)

		// Continue to the next middleware or handler in the chain.
		ginContext.Next()
	}
}
