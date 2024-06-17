package middleware

import (
	"context"

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
		// Create a context with timeout.
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the refresh token from the request headers or cookies.
		refreshToken, tokenError := extractRefreshToken(ginContext)
		if validator.IsError(tokenError) {
			// Abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, tokenError, constants.StatusUnauthorized)
			return
		}

		// Get the refresh token configuration.
		applicationConfig := config.AppConfig.RefreshToken

		// Validate the JWT token.
		userTokenPayload := domainUtility.ValidateJWTToken(location+".RefreshTokenAuthenticationMiddleware", refreshToken, applicationConfig.PublicKey)
		if validator.IsError(userTokenPayload.Error) {
			// Handle token validation error and respond with an unauthorized status and JSON error.
			httpAuthorizationError := httpError.NewHttpAuthorizationErrorView("", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, httpAuthorizationError, constants.StatusUnauthorized)
			return
		}

		// Store user information in the context.
		ctx = context.WithValue(ctx, constants.UserIDContext, userTokenPayload.Data.UserID)
		ctx = context.WithValue(ctx, constants.UserRoleContext, userTokenPayload.Data.Role)

		// Update the request's context with the new context containing user information.
		ginContext.Request = ginContext.Request.WithContext(ctx)

		// Continue to the next middleware or handler in the chain.
		ginContext.Next()
	}
}
