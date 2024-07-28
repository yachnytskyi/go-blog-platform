package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userView "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
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

// UserContextMiddleware is a middleware for retrieving user information based on the user ID from the context.
// This middleware extracts the user ID, fetches the corresponding user details, and stores them in the request context.
// Note: This middleware should be placed after the AuthenticationMiddleware in the middleware chain to ensure the user ID is available.
func UserContextMiddleware(logger interfaces.Logger, userUseCase user.UserUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		userID, _ := ctx.Value(constants.ID).(string)
		if userID == "" {
			httpInternalError := httpError.NewHTTPInternalError(location+"UserContextMiddleware.ctx.Value", constants.IDContextMissing)
			abortWithStatusJSON(ginContext, logger, httpInternalError, constants.StatusUnauthorized)
			return
		}

		user := userUseCase.GetUserById(ctx, userID)
		if validator.IsError(user.Error) {
			abortWithStatusJSON(ginContext, logger, user.Error, constants.StatusUnauthorized)
			return
		}

		userView := userView.UserToUserViewMapper(user.Data)
		ctx = context.WithValue(ctx, constants.User, userView)
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}
