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
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func RefreshTokenAuthenticationMiddleware(userUseCase user.UserUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Check context timeout.
		contextError := commonUtility.HandleWithContextError(location+"RefreshTokenAuthenticationMiddleware", ctx)
		if validator.IsError(contextError) {
			abortWithStatusJSON(ginContext, contextError, http.StatusUnauthorized)
		}

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
		userID, validateRefreshTokenError := domainUtility.ValidateJWTToken(refreshToken, applicationConfig.PublicKey)
		if validator.IsError(validateRefreshTokenError) {
			// Handle token validation error and respond with an unauthorized status and JSON error.
			httpAuthorizationError := httpError.NewHttpAuthorizationErrorView(constants.EmptyString, constants.LoggingErrorNotification)
			abortWithStatusJSON(ginContext, httpAuthorizationError, http.StatusUnauthorized)
			return
		}

		// Get the user information from the user use case.
		user := userUseCase.GetUserById(ctx, fmt.Sprint(userID))
		if validator.IsError(user.Error) {
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
