package utility

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"

	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
)

const (
	accessToToken = "access_token"
	authorization = "Authorization"
	bearer        = "Bearer"
)

func DeserializeUser(userUseCase user.UseCase) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var accessToken string
		cookie, cookieError := ctx.Cookie(accessToToken)
		authorizationHeader := ctx.Request.Header.Get(authorization)
		fields := strings.Fields(authorizationHeader)
		if validator.IsSliceNotEmpty(fields) && validator.CheckMatchStrings(fields[0], bearer) {
			accessToken = fields[1]
		} else if validator.IsValueNil(cookieError) {
			accessToken = cookie
		}
		if validator.IsStringEmpty(accessToken) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		config, _ := config.LoadConfig(".")
		userID, validateTokenError := httpUtility.ValidateToken(accessToken, config.AccessTokenPublicKey)
		if validator.IsValueNotNil(validateTokenError) {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": validateTokenError.Error()})
			return
		}

		context := ctx.Request.Context()
		user, getUSerByIdError := userUseCase.GetUserById(context, fmt.Sprint(userID))
		if validator.IsValueNotNil(getUSerByIdError) {

			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "You are not logged in"})
			return
		}

		userMappedToUserView := userViewModel.UserToUserViewMapper(user)
		ctx.Set("currentUser", userMappedToUserView)
		ctx.Next()
	}
}
