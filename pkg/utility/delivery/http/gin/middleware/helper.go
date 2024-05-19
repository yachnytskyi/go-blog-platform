package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// extractAccessToken extracts the access token from the request headers or cookies.
// Parameters:
// - ginContext: The Gin context containing the HTTP request.
// Returns:
// - The extracted access token as a string.
// - An error if the access token is not found or invalid.
func extractAccessToken(ginContext *gin.Context) (string, error) {
	// Attempt to retrieve the access token from the cookie.
	cookie, cookieError := ginContext.Cookie(constants.AccessTokenValue)
	// Retrieve the Authorization header from the request.
	authorizationHeader := ginContext.Request.Header.Get(authorization)
	fields := strings.Fields(authorizationHeader)

	// Check if the Authorization header contains a Bearer token.
	if validator.IsSliceNotEmpty(fields) && fields[0] == bearer {
		// If a Bearer token is present, set the access token.
		return fields[1], nil
	} else if cookieError == nil {
		// If no Bearer token in the Authorization header, try to get the token from the cookie.
		return cookie, nil
	}

	// If access token is still empty, create and log an HTTP authorization error.
	httpAuthorizationError := httpError.NewHttpAuthorizationErrorView(location+"extractAccessToken.accessToken", constants.LoggingErrorNotification)
	logging.Logger(httpAuthorizationError)
	return "", httpAuthorizationError
}

// extractRefreshToken extracts the refresh token from the request cookies.
// Parameters:
// - ginContext: The Gin context containing the HTTP request.
// Returns:
// - The extracted refresh token as a string.
// - An error if the refresh token is not found or invalid.
func extractRefreshToken(ginContext *gin.Context) (string, error) {
	// Attempt to retrieve the refresh token from the cookie.
	refreshToken, refreshTokenError := ginContext.Cookie(constants.RefreshTokenValue)
	if validator.IsError(refreshTokenError) {
		// If refresh token is missing, create and log an HTTP authorization error.
		httpAuthorizationError := httpError.NewHttpAuthorizationErrorView(location+"extractRefreshToken.refreshToken", constants.LoggingErrorNotification)
		logging.Logger(httpAuthorizationError)
		return "", httpAuthorizationError
	}

	// Return the extracted refresh token.
	return refreshToken, nil
}

// abortWithStatusJSON aborts the request, logs the error, and responds with a JSON error.
// Parameters:
// - ginContext: The Gin context used to generate the JSON response.
// - err: The error to be included in the response.
// - httpCode: The HTTP status code to be set in the response.
func abortWithStatusJSON(ginContext *gin.Context, err error, httpCode int) {
	logging.Logger(err)
	jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(err))
	ginContext.AbortWithStatusJSON(httpCode, jsonResponse)
}
