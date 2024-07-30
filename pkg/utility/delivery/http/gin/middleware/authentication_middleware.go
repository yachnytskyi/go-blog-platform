package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// AuthenticationMiddleware is a Gin middleware for handling user authentication using JWT tokens.
func AuthenticationMiddleware(configInstance interfaces.Config, logger interfaces.Logger) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the access token from the request headers or cookies.
		accessToken := extractAccessToken(ginContext, location+"AuthenticationMiddleware")
		if validator.IsError(accessToken.Error) {
			abortWithStatusJSON(ginContext, logger, accessToken.Error, constants.StatusUnauthorized)
			return
		}

		// Extract the access token from the request headers or cookies.
		config := configInstance.GetConfig()
		userTokenPayload := utility.ValidateJWTToken(logger, location+"AuthenticationMiddleware", accessToken.Data, config.AccessToken.PublicKey)
		if validator.IsError(userTokenPayload.Error) {
			httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+"AuthenticationMiddleware.ValidateJWTToken", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, logger, httpAuthorizationError, constants.StatusUnauthorized)
			return
		}

		ctx = context.WithValue(ctx, constants.ID, userTokenPayload.Data.UserID)
		ctx = context.WithValue(ctx, constants.UserRole, userTokenPayload.Data.Role)
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}
