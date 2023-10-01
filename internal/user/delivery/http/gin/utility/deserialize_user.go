package utility

import (
	"fmt"
	"net/http"
	"strings"

	gin "github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	accessToToken = "access_token"
	authorization = "Authorization"
	bearer        = "Bearer"
)

func DeserializeUser(userUseCase user.UserUseCase) gin.HandlerFunc {
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

		context := ctx.Request.Context()
		user, getUSerByIdError := userUseCase.GetUserById(context, fmt.Sprint(userID))
		if validator.IsErrorNotNil(getUSerByIdError) {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		userMappedToUserView := userViewModel.UserToUserViewMapper(user)
		ctx.Set("currentUser", userMappedToUserView)
		ctx.Next()
	}
}
