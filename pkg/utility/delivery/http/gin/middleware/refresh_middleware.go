package middleware

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func RefreshTokenAuthenticationMiddleware(userUseCase user.UserUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Extract the refresh token from the request.
		refreshToken, tokenError := extractRefreshToken(ginContext)
		if validator.IsErrorNotNil(tokenError) {
			// Abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, tokenError, http.StatusUnauthorized)
			return
		}

		// Get the application configuration.
		applicationConfig := config.AppConfig

		// Validate the JWT token.
		userID, validateRefreshTokenError := domainUtility.ValidateJWTToken(refreshToken, applicationConfig.RefreshToken.PublicKey)
		if validator.IsErrorNotNil(validateRefreshTokenError) {
			// Handle token validation error and respond with an unauthorized status and JSON error.
			httpAuthorizationError := httpError.NewHttpAuthorizationErrorView(constants.EmptyString, constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, httpAuthorizationError, http.StatusUnauthorized)
			return
		}

		// Check for a deadline error using the handleDeadlineExceeded function.
		// If a deadline error occurred, respond with a timeout status.
		deadlineError := handleDeadlineExceeded(ctx)
		if validator.IsValueNotNil(deadlineError) {
			// Use the abortWithStatusJSON function to handle the deadline error by sending
			// a JSON response with an appropriate HTTP status code.
			abortWithStatusJSON(ginContext, deadlineError, http.StatusUnauthorized)
		}

		// Get the user information from the user use case.
		user := userUseCase.GetUserById(ctx, fmt.Sprint(userID))
		if validator.IsErrorNotNil(user.Error) {
			// Handle user retrieval error and respond with an unauthorized status and JSON error.
			handledError := httpError.HandleError(user.Error)
			abortWithStatusJSON(ginContext, handledError, http.StatusUnauthorized)
			return
		}

		// Set user-related information in the Gin context for downstream handlers.
		ginContext.Set(constants.UserIDContext, userID)
		ginContext.Set(constants.UserContext, userViewModel.UserToUserViewMapper(user.Data))

		// Continue to the next middleware or handler in the chain.
		ginContext.Next()
	}
}
