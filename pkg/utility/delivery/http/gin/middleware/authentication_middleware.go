package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location             = "pkg.utility.delivery.http.gin.middleware."
	authorization        = "Authorization"
	bearer               = "Bearer"
	userIDContextMissing = "User ID context value is missing or empty."
)

// AuthenticationMiddleware is a Gin middleware for handling user authentication using JWT tokens.
// Returns a Gin middleware handler function.
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Create a context with timeout.
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the refresh token from the request headers or cookies.
		accessToken := extractAccessToken(ginContext)
		if validator.IsError(accessToken.Error) {
			// Abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, accessToken.Error, constants.StatusUnauthorized)
			return
		}

		// Get the application configuration.
		accessTokenConfig := config.GetAccessConfig()
		// Validate the JWT token.
		userTokenPayload := domainUtility.ValidateJWTToken(location+".AuthenticationMiddleware", accessToken.Data, accessTokenConfig.PublicKey)
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

// UserContextMiddleware is a middleware for retrieving user information based on the user ID from the context.
// This middleware extracts the user ID, fetches the corresponding user details, and stores them in the request context.
// Note: This middleware should be placed after the AuthenticationMiddleware in the middleware chain to ensure the user ID is available.
func UserContextMiddleware(userUseCase user.UserUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Create a context with a timeout.
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Retrieve the user ID from the context.
		userID, _ := ctx.Value(constants.UserIDContext).(string)
		if userID == "" {
			httpInternalErrorView := httpError.NewHttpInternalErrorView(location+"UserContextMiddleware.ctx.Value", userIDContextMissing)
			abortWithStatusJSON(ginContext, httpInternalErrorView, constants.StatusUnauthorized)
			return
		}

		// Fetch user details using the user use case.
		user := userUseCase.GetUserById(ctx, userID)
		if validator.IsError(user.Error) {
			// If an error occurs during user retrieval, abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, user.Error, constants.StatusUnauthorized)
			return
		}

		// Map user data to a user view model.
		userView := userViewModel.UserToUserViewMapper(user.Data)

		// Store user information in the context.
		ctx = context.WithValue(ctx, constants.UserContext, userView)

		// Update the request's context with the new context containing user information.
		ginContext.Request = ginContext.Request.WithContext(ctx)

		// Continue to the next middleware or handler.
		ginContext.Next()
	}
}
