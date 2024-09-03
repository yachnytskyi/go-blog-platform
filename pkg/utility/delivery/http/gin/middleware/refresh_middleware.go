package middleware

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// RefreshTokenMiddleware is a Gin middleware for handling user authentication using refresh tokens.
func RefreshTokenMiddleware(config *config.ApplicationConfig, logger interfaces.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the refresh token from the request headers or cookies.
		refreshToken := extractToken(ginContext, location+"RefreshTokenMiddleware", constants.RefreshTokenValue)
		if validator.IsError(refreshToken.Error) {
			abortWithStatusJSON(ginContext, logger, refreshToken.Error, http.StatusUnauthorized)
			return
		}

		// Validate the JWT token using the public key from the configuration.
		userTokenPayload := utility.ValidateJWTToken(
			logger,
			location+"RefreshTokenMiddleware",
			refreshToken.Data,
			config.RefreshToken.PublicKey,
		)
		if validator.IsError(userTokenPayload.Error) {
			httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+"RefreshTokenMiddleware.ValidateJWTToken", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, logger, httpAuthorizationError, http.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, constants.ID, userTokenPayload.Data.UserID)
		ctx = context.WithValue(ctx, constants.UserRole, userTokenPayload.Data.Role)
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}
