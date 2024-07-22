package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// RefreshTokenAuthenticationMiddleware is a Gin middleware for handling user authentication using refresh tokens.
func RefreshTokenAuthenticationMiddleware(logger applicationModel.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the refresh token from the request headers or cookies.
		refreshToken := extractRefreshToken(ginContext, logger)
		if validator.IsError(refreshToken.Error) {
			abortWithStatusJSON(ginContext, logger, refreshToken.Error, constants.StatusUnauthorized)
			return
		}

		// Extract the refresh token from the request headers or cookies.
		refreshTokenConfig := config.GetRefreshConfig()
		userTokenPayload := domainUtility.ValidateJWTToken(logger, location+".RefreshTokenAuthenticationMiddleware", refreshToken.Data, refreshTokenConfig.PublicKey)
		if validator.IsError(userTokenPayload.Error) {
			httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+"RefreshTokenAuthenticationMiddleware.ValidateJWTToken", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, logger, httpAuthorizationError, constants.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, constants.ID, userTokenPayload.Data.UserID)
		ctx = context.WithValue(ctx, constants.UserRole, userTokenPayload.Data.Role)
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}
