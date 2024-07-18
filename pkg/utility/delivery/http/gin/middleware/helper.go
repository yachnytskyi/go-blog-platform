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
func extractAccessToken(ginContext *gin.Context, location string) commonModel.Result[string] {
	cookie, cookieError := ginContext.Cookie(constants.AccessTokenValue)
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
	httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+".extractAccessToken.accessToken", constants.LoggingErrorNotification)
	logging.Logger(httpAuthorizationError)
	return commonModel.NewResultOnFailure[string](httpAuthorizationError)
}

// extractRefreshToken extracts the refresh token from the request cookies.
func extractRefreshToken(ginContext *gin.Context) commonModel.Result[string] {
	refreshToken, refreshTokenError := ginContext.Cookie(constants.RefreshTokenValue)
	if validator.IsError(refreshTokenError) {
		httpAuthorizationError := httpError.NewHTTPAuthorizationError(location+".extractRefreshToken.refreshToken", constants.LoggingErrorNotification)
		logging.Logger(httpAuthorizationError)
		return commonModel.NewResultOnFailure[string](httpAuthorizationError)
	}

	return commonModel.NewResultOnSuccess[string](refreshToken)
}

// abortWithStatusJSON aborts the request, logs the error, and responds with a JSON error.
func abortWithStatusJSON(ginContext *gin.Context, err error, httpCode int) {
	logging.Logger(err)
	jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(err))
	ginContext.AbortWithStatusJSON(httpCode, jsonResponse)
}
