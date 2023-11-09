package middleware

import (
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
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	authorization = "Authorization"
	bearer        = "Bearer"
	firstElement  = 0
	nextElement   = 1
)

// AuthMiddleware is a Gin middleware for handling user authentication using JWT tokens.
func AuthMiddleware(userUseCase user.UserUseCase) gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		// Extract the access token from the request.
		accessToken, extractAccessTokenError := extractAccessToken(ginContext)
		if extractAccessTokenError != nil {
			// Abort the request with an unauthorized status and respond with a JSON error.
			abortWithStatusJSON(ginContext, extractAccessTokenError, http.StatusUnauthorized)
			return
		}

		// Get the application configuration.
		applicationConfig := config.AppConfig

		// Validate the JWT token.
		userID, validateTokenError := domainUtility.ValidateJWTToken(accessToken, applicationConfig.AccessToken.PublicKey)
		if validator.IsErrorNotNil(validateTokenError) {
			// Handle token validation error and respond with an unauthorized status and JSON error.
			err := httpError.HandleError(validateTokenError)
			abortWithStatusJSON(ginContext, err, http.StatusUnauthorized)
			return
		}

		// Get the user information from the user use case.
		context := ginContext.Request.Context()
		user := userUseCase.GetUserById(context, fmt.Sprint(userID))
		if validator.IsErrorNotNil(user.Error) {
			// Handle user retrieval error and respond with an unauthorized status and JSON error.
			err := httpError.HandleError(user.Error)
			abortWithStatusJSON(ginContext, err, http.StatusUnauthorized)
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
	accessToken := ""
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
	if validator.IsStringEmpty(accessToken) {
		// If access token is empty, create and log an HTTP authorization error.
		httpAuthorizationError := httpError.NewHttpAuthorizationErrorView(constants.LoggingErrorNotification)
		logging.Logger(httpAuthorizationError)
		return "", httpAuthorizationError
	}

	// Return the extracted access token.
	return accessToken, nil
}

// abortWithStatusJSON aborts the request, logs the error, and responds with a JSON error.
func abortWithStatusJSON(ginContext *gin.Context, err any, httpCode int) {
	logging.Logger(err)
	jsonResponse := httpModel.NewJsonResponseOnFailure(err)
	ginContext.AbortWithStatusJSON(httpCode, jsonResponse)
}
