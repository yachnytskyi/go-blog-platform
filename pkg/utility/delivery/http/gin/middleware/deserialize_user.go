package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	authorization = "Authorization"
	bearer        = "Bearer"
)

func DeserializeUser(userUseCase user.UserUseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, cookieError := ctx.Cookie(constant.AccessTokenValue)
		authorizationHeader := ctx.Request.Header.Get(authorization)
		fields := strings.Fields(authorizationHeader)
		if validator.IsSliceNotEmpty(fields) && fields[0] == bearer {
			accessToken = fields[1]
		} else if validator.IsErrorNil(cookieError) {
			accessToken = cookie
		}
		if validator.IsStringEmpty(accessToken) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		applicationConfig := config.AppConfig
		userID, err := httpUtility.ValidateToken(accessToken, applicationConfig.AccessToken.PublicKey)
		if validator.IsErrorNotNil(err) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		context := ctx.Request.Context()
		user := userUseCase.GetUserById(context, fmt.Sprint(userID))
		if validator.IsErrorNotNil(user.Error) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}
		userMappedToUserView := userViewModel.UserToUserViewMapper(user.Data)
		ctx.Set(constant.UserIDContext, userID)
		ctx.Set(constant.UserContext, userMappedToUserView)
		ctx.Next()
	}
}
