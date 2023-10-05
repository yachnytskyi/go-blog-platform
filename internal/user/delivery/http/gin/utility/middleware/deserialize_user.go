package middleware

import (
	"net/http"
	"strings"

	gin "github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	accessToToken = "access_token"
	authorization = "Authorization"
	bearer        = "Bearer"
)

func DeserializeUser() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, cookieError := ctx.Cookie(accessToToken)
		authorizationHeader := ctx.Request.Header.Get(authorization)
		fields := strings.Fields(authorizationHeader)
		if validator.IsSliceNotEmpty(fields) && validator.AreStringsEqual(fields[0], bearer) {
			accessToken = fields[1]
		} else if validator.IsErrorNil(cookieError) {
			accessToken = cookie
		}
		if validator.IsStringEmpty(accessToken) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		applicationConfig := config.AppConfig
		userID, validateTokenError := httpUtility.ValidateToken(accessToken, applicationConfig.AccessToken.PublicKey)
		if validator.IsErrorNotNil(validateTokenError) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": validateTokenError.Error()})
			return
		}
		ctx.Set("userID", userID)
		ctx.Next()
	}
}
