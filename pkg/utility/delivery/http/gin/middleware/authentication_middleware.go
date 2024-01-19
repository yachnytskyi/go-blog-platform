package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location      = "pkg.utility.delivery.http.gin.middleware."
	authorization = "Authorization"
	bearer        = "Bearer"
	firstElement  = 0
	nextElement   = 1
)

// AuthMiddleware is a Gin middleware for handling user authentication using JWT tokens.
func AuthenticationMiddleware(userUseCase user.UserUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
		defer cancel()

		// Check context timeout.
		contextError := commonUtility.HandleWithContextError(location+"AuthMiddleware", ctx)
		if validator.IsError(contextError) {
			abortWithStatusJSON(ginContext, contextError, http.StatusUnauthorized)
		}

		// Extract the access token from the request.
		accessToken, tokenError := extractAccessToken(ginContext)
		if validator.IsError(tokenError) {
			// Abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, tokenError, http.StatusUnauthorized)
			return
		}

		// Get the application configuration.
		accessTokenConfig := config.AppConfig.AccessToken

		// Validate the JWT token.
		userID, validateAccessTokenError := domainUtility.ValidateJWTToken(accessToken, accessTokenConfig.PublicKey)
		if validator.IsError(validateAccessTokenError) {
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

// extractAccessToken extracts the access token from the request headers or cookies.
func extractAccessToken(ginContext *gin.Context) (string, error) {
	// Initialize variables to store the access token.
	// Attempt to retrieve the access token from the cookie.
	// Retrieve the Authorization header from the request.
	var accessToken string
	cookie, cookieError := ginContext.Cookie(constants.AccessTokenValue)
	authorizationHeader := ginContext.Request.Header.Get(authorization)
	fields := strings.Fields(authorizationHeader)

	// Check if the Authorization header contains a Bearer token.
	if validator.IsSliceNotEmpty(fields) && fields[firstElement] == bearer {
		// If a Bearer token is present, set the access token.
		accessToken = fields[nextElement]
	} else if cookieError == nil {
		// If no Bearer token in the Authorization header, try to get the token from the cookie.
		accessToken = cookie
	}

	// Check if the access token is still empty.
	if accessToken == constants.EmptyString {
		// If access token is empty, create and log an HTTP authorization error.
		httpAuthorizationError := httpError.NewHttpAuthorizationErrorView(location+"extractAcessToken.accessToken", constants.LoggingErrorNotification)
		logging.Logger(httpAuthorizationError)
		return constants.EmptyString, httpAuthorizationError
	}

	// Return the extracted access token.
	return accessToken, nil
}

// extractRefreshToken extracts the refresh token from the request cookies.
func extractRefreshToken(ginContext *gin.Context) (string, error) {
	// Attempt to retrieve the refresh token from the cookie.
	refreshToken, refreshTokenError := ginContext.Cookie(constants.RefreshTokenValue)
	if validator.IsError(refreshTokenError) {
		// If refresh token is missing, create and log an HTTP authorization error.
		httpAuthorizationError := httpError.NewHttpAuthorizationErrorView(location+"extractRefreshToken.refreshToken", constants.LoggingErrorNotification)
		logging.Logger(httpAuthorizationError)
		return constants.EmptyString, httpAuthorizationError
	}

	// Return the extracted refresh token.
	return refreshToken, nil
}

// abortWithStatusJSON aborts the request, logs the error, and responds with a JSON error.
func abortWithStatusJSON(ginContext *gin.Context, err error, httpCode int) {
	logging.Logger(err)
	jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(err))
	ginContext.AbortWithStatusJSON(httpCode, jsonResponse)
}
