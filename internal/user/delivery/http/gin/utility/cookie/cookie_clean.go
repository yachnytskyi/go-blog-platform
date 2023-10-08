package cookie

import (
	"github.com/gin-gonic/gin"
	constant "github.com/yachnytskyi/golang-mongo-grpc/config/constant"
)

func CleanCookies(ctx *gin.Context) {
	ctx.SetCookie(constant.AccessTokenValue, "", constant.LogoutMaxAgeValue, "/", constant.TokenDomainValue, false, true)
	ctx.SetCookie("refresh_token", "", constant.LogoutMaxAgeValue, "/", constant.TokenDomainValue, false, true)
	ctx.SetCookie(constant.LoggedInValue, "", constant.LogoutMaxAgeValue, "/", constant.TokenDomainValue, false, true)
}
