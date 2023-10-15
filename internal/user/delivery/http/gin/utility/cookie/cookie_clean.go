package cookie

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

func CleanCookies(ctx *gin.Context) {
	ctx.SetCookie(constants.AccessTokenValue, "", constants.LogoutMaxAgeValue, "/", constants.TokenDomainValue, false, true)
	ctx.SetCookie("refresh_token", "", constants.LogoutMaxAgeValue, "/", constants.TokenDomainValue, false, true)
	ctx.SetCookie(constants.LoggedInValue, "", constants.LogoutMaxAgeValue, "/", constants.TokenDomainValue, false, true)
}
