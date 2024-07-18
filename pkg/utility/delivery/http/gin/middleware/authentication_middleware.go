package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userView "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "pkg.utility.delivery.http.gin.middleware."
)

// AuthenticationMiddleware is a Gin middleware for handling user authentication using JWT tokens.
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the access token from the request headers or cookies.
		accessToken := extractAccessToken(ginContext, location)
		if validator.IsError(accessToken.Error) {
			abortWithStatusJSON(ginContext, accessToken.Error, constants.StatusUnauthorized)
			return
		}

		// Extract the access token from the request headers or cookies.
		accessTokenConfig := config.GetAccessConfig()
		userTokenPayload := domainUtility.ValidateJWTToken(location+"AuthenticationMiddleware", accessToken.Data, accessTokenConfig.PublicKey)
		if validator.IsError(userTokenPayload.Error) {
			httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+"AuthenticationMiddleware.ValidateJWTToken", constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, httpAuthorizationError, constants.StatusUnauthorized)
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
func UserContextMiddleware(userUseCase user.UserUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		userID, _ := ctx.Value(constants.ID).(string)
		if len(userID) == 0 {
			httpInternalError := httpError.NewHTTPInternalError(location+"UserContextMiddleware.ctx.Value", constants.IDContextMissing)
			abortWithStatusJSON(ginContext, httpInternalError, constants.StatusUnauthorized)
			return
		}

		user := userUseCase.GetUserById(ctx, userID)
		if validator.IsError(user.Error) {
			abortWithStatusJSON(ginContext, user.Error, constants.StatusUnauthorized)
			return
		}

		userView := userView.UserToUserViewMapper(user.Data)
		ctx = context.WithValue(ctx, constants.User, userView)
		ginContext.Request = ginContext.Request.WithContext(ctx)
		ginContext.Next()
	}
}
