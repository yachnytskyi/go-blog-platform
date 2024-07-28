package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "pkg.utility.delivery.http.gin.middleware."
)

// extractAccessToken extracts the access token from the request headers or cookies.
func extractAccessToken(ginContext *gin.Context, location string) common.Result[string] {
	cookie, cookieError := ginContext.Cookie(constants.AccessTokenValue)
	authorizationHeader := ginContext.Request.Header.Get(constants.Authorization)
	fields := strings.Fields(authorizationHeader)

	// Check if the Authorization header contains a Bearer token.
	if len(fields) > 0 && fields[0] == constants.Bearer {
		// If a Bearer token is present, set the access token.
		return common.NewResultOnSuccess[string](fields[1])
	} else if cookieError == nil {
		// If no Bearer token in the Authorization header, try to get the token from the cookie.
		return common.NewResultOnSuccess[string](cookie)
	}

	// If access token is still empty, create and log an HTTP authorization error.
	return common.NewResultOnFailure[string](httpError.NewHTTPAuthorizationError(location+".extractAccessToken.accessToken", constants.LoggingErrorNotification))
}

// extractRefreshToken extracts the refresh token from the request cookies.
func extractRefreshToken(ginContext *gin.Context, location string) common.Result[string] {
	refreshToken, refreshTokenError := ginContext.Cookie(constants.RefreshTokenValue)
	if validator.IsError(refreshTokenError) {
		return common.NewResultOnFailure[string](httpError.NewHTTPAuthorizationError(location+".extractRefreshToken.refreshToken", constants.LoggingErrorNotification))
	}

	return common.NewResultOnSuccess[string](refreshToken)
}

// abortWithStatusJSON aborts the request, logs the error, and responds with a JSON error.
func abortWithStatusJSON(ginContext *gin.Context, logger interfaces.Logger, err error, httpCode int) {
	logger.Error(err)
	jsonResponse := httpModel.NewJSONResponseOnFailure(httpError.HandleError(err))
	ginContext.AbortWithStatusJSON(httpCode, jsonResponse)
}
