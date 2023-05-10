package utility

import "github.com/gin-gonic/gin"

func LoginSetCookies(ctx *gin.Context, accessToken string, accessTokenMaxAge int, refreshToken string, refreshTokenMaxAge int) {
	ctx.SetCookie("access_token", accessToken, accessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, refreshTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", accessTokenMaxAge, "/", "localhost", false, false)
}

func ResetUserPasswordSetCookies(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
}

func RefreshAccessTokenSetCookies(ctx *gin.Context, accessToken string, accessTokenMaxAge int) {
	ctx.SetCookie("access_token", accessToken, accessTokenMaxAge, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", accessTokenMaxAge, "/", "localhost", false, false)
}

func LogoutSetCookies(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)
}
