package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	authorization = "Authorization"
	bearer        = "Bearer"
)

func AuthContextMiddleware(userUseCase user.UserUseCase) gin.HandlerFunc {
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
		if validator.IsStringEmpty(accessToken) {
			authorizationError := httpError.NewHttpAuthorizationErrorView(constants.LoggingErrorNotification)
			logging.Logger(authorizationError)
			jsonResponse := httpModel.NewJsonResponseOnFailure(authorizationError)
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, jsonResponse)
			return
		}

		applicationConfig := config.AppConfig
		userID, validateTokenError := domainUtility.ValidateJWTToken(accessToken, applicationConfig.AccessToken.PublicKey)
		if validator.IsErrorNotNil(validateTokenError) {
			jsonResponse := httpModel.NewJsonResponseOnFailure(validateTokenError)
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, jsonResponse)
			return
		}
		context := ginContext.Request.Context()
		user := userUseCase.GetUserById(context, fmt.Sprint(userID))
		if validator.IsErrorNotNil(user.Error) {
			jsonResponse := httpModel.NewJsonResponseOnFailure(user.Error)
			ginContext.AbortWithStatusJSON(http.StatusUnauthorized, jsonResponse)
			return
		}
		ginContext.Set(constants.UserIDContext, userID)
		ginContext.Set(constants.UserContext, userViewModel.UserToUserViewMapper(user.Data))
		ginContext.Next()
	}
}
