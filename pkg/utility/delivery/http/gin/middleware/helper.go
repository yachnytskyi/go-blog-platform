package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	http "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
)

const (
	location = "pkg.utility.delivery.http.gin.middleware."
)

// extractToken extracts the token from the request headers or cookies.
func extractToken(ginContext *gin.Context, location, tokenType string) common.Result[string] {
	// Check if the cookie contains the token.
	cookie, cookieError := ginContext.Cookie(tokenType)
	if cookieError == nil {
		return common.NewResultOnSuccess[string](cookie)
	}

	// If no token in the cookie, try to get the token from the Authorization header.
	authorizationHeader := ginContext.Request.Header.Get(constants.Authorization)
	if strings.HasPrefix(authorizationHeader, constants.Bearer) {
		authorizationHeader = authorizationHeader[len(constants.Bearer):]
		if len(authorizationHeader) > 0 {
			return common.NewResultOnSuccess[string](authorizationHeader)
		}
	}

	// If no token was found, return a failure with an HTTP authorization error.
	return common.NewResultOnFailure[string](delivery.NewHTTPAuthorizationError(location+".extractToken.token", constants.LoggingErrorNotification))
}

// abortWithStatusJSON aborts the request, logs the error, and responds with a JSON error.
func abortWithStatusJSON(ginContext *gin.Context, logger interfaces.Logger, err error, httpCode int) {
	logger.Error(err)
	jsonResponse := http.NewJSONResponseOnFailure(delivery.HandleError(err))
	ginContext.AbortWithStatusJSON(httpCode, jsonResponse)
}
