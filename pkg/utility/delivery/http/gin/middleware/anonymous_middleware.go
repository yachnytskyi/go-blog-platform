package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

func AnonymousContextMiddleware() gin.HandlerFunc {
	return func(ginContext *gin.Context) {
		var accessToken string
		cookie, cookieError := ginContext.Cookie(constant.AccessTokenValue)
		authorizationHeader := ginContext.Request.Header.Get(authorization)
		fields := strings.Fields(authorizationHeader)
		if validator.IsSliceNotEmpty(fields) && fields[0] == bearer {
			accessToken = fields[1]
		} else if validator.IsErrorNil(cookieError) {
			accessToken = cookie
		}
		if validator.IsStringNotEmpty(accessToken) {
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are already registered, and registration is not allowed for existing users."})
		}
		ginContext.Next()
	}
}
