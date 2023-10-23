package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func AnonymousContextMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		var accessToken string
		cookie, cookieError := ginContext.Cookie(constants.AccessTokenValue)
		authorizationHeader := ginContext.Request.Header.Get(authorization)
		fields := strings.Fields(authorizationHeader)
		if validator.IsSliceNotEmpty(fields) && fields[0] == bearer {
			accessToken = fields[1]
		} else if validator.IsErrorNil(cookieError) {
			accessToken = cookie
		}
		if validator.IsStringNotEmpty(accessToken) {
			authorizationError := httpError.NewHttpAuthorizationErrorView(constants.AlreadyRegisteredNotification)
			logging.Logger(authorizationError)
			jsonResponse := httpModel.NewJsonResponseOnFailure(authorizationError)
			httpModel.SetStatus(&jsonResponse)
			ginContext.AbortWithStatusJSON(http.StatusForbidden, jsonResponse)
			return
		}
		ginContext.Next()
	}
}
