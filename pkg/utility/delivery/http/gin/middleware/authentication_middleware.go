package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/interfaces"
	utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// AuthenticationMiddleware is a Gin middleware for handling user authentication using JWT tokens.
func AuthenticationMiddleware(config *config.ApplicationConfig, logger interfaces.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the access token from the request headers or cookies.
		accessToken := extractToken(ginContext, location+"AuthenticationMiddleware", constants.AccessTokenValue)
		if validator.IsError(accessToken.Error) {
			abortWithStatusJSON(ginContext, logger, accessToken.Error, http.StatusUnauthorized)
			return
		}

		// Validate the JWT token using the public key from the configuration.
		userTokenPayload := utility.ValidateJWTToken(
			logger,
			location+"AuthenticationMiddleware",
			accessToken.Data,
			config.AccessToken.PublicKey,
		)
		if validator.IsError(userTokenPayload.Error) {
			httpAuthorizationError := delivery.NewHTTPAuthorizationError(location+"AuthenticationMiddleware.ValidateJWTToken", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, logger, httpAuthorizationError, http.StatusUnauthorized)
			return
		}

		// Store the user's ID and role in the request context.
		ctx = context.WithValue(ctx, constants.ID, userTokenPayload.Data.UserID)
		ctx = context.WithValue(ctx, constants.UserRole, userTokenPayload.Data.Role)
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}
