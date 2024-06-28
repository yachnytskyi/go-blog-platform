package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

// extractAccessToken extracts the access token from the request headers or cookies.
//
// Parameters:
// - ginContext: The Gin context containing the HTTP request.
//
// Returns:
// - A Result containing the extracted access token as a string.
// - A Result containing an error if the access token is not found or invalid.
func extractAccessToken(ginContext *gin.Context) commonModel.Result[string] {
	// Attempt to retrieve the access token from the cookie.
	cookie, cookieError := ginContext.Cookie(constants.AccessTokenValue)
	// Retrieve the Authorization header from the request.
	authorizationHeader := ginContext.Request.Header.Get(constants.Authorization)
	fields := strings.Fields(authorizationHeader)

	// Check if the Authorization header contains a Bearer token.
	if validator.IsSliceNotEmpty(fields) && fields[0] == constants.Bearer {
		// If a Bearer token is present, set the access token.
		return commonModel.NewResultOnSuccess[string](fields[1])
	} else if cookieError == nil {
		// If no Bearer token in the Authorization header, try to get the token from the cookie.
		return commonModel.NewResultOnSuccess[string](cookie)
	}

	// If access token is still empty, create and log an HTTP authorization error.
	httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+"extractAccessToken.accessToken", constants.LoggingErrorNotification)
	logging.Logger(httpAuthorizationError)
	return commonModel.NewResultOnFailure[string](httpAuthorizationError)
}

// extractRefreshToken extracts the refresh token from the request cookies.
//
// Parameters:
// - ginContext: The Gin context containing the HTTP request.
//
// Returns:
// - A Result containing the extracted refresh token as a string.
// - A Result containing an error if the refresh token is not found or invalid.
func extractRefreshToken(ginContext *gin.Context) commonModel.Result[string] {
	// Attempt to retrieve the refresh token from the cookie.
	refreshToken, refreshTokenError := ginContext.Cookie(constants.RefreshTokenValue)
	if validator.IsError(refreshTokenError) {
		// If refresh token is missing, create and log an HTTP authorization error.
		httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+"extractRefreshToken.refreshToken", constants.LoggingErrorNotification)
		logging.Logger(httpAuthorizationError)
		return commonModel.NewResultOnFailure[string](httpAuthorizationError)
	}

	// Return the extracted refresh token.
	return commonModel.NewResultOnSuccess[string](refreshToken)
}

// abortWithStatusJSON aborts the request, logs the error, and responds with a JSON error.
//
// Parameters:
// - ginContext: The Gin context used to generate the JSON response.
// - err: The error to be included in the response.
// - httpCode: The HTTP status code to be set in the response.
func abortWithStatusJSON(ginContext *gin.Context, err error, httpCode int) {
	logging.Logger(err)
	jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(err))
	ginContext.AbortWithStatusJSON(httpCode, jsonResponse)
}
