package cookie

import (
	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

// CleanCookies is a helper function for Gin controllers that clears specific cookies.
func CleanCookies(ctx *gin.Context, path string) {
	// Retrieve application configuration for cookie settings.
	securityConfig := config.GetSecurityConfig()

	// Clear the access token cookie by setting its value to an empty string and
	// configuring it with the logout max age, path, and domain.
	ctx.SetCookie(constants.AccessTokenValue, "", constants.LogoutMaxAgeValue, path, constants.TokenDomainValue, securityConfig.CookieSecure, securityConfig.HTTPOnly)

	// Clear the refresh token cookie in a similar manner.
	ctx.SetCookie(constants.RefreshTokenValue, "", constants.LogoutMaxAgeValue, path, constants.TokenDomainValue, securityConfig.CookieSecure, securityConfig.HTTPOnly)

	// Clear the "loggedIn" cookie in the same way.
	ctx.SetCookie(constants.LoggedInValue, "", constants.LogoutMaxAgeValue, path, constants.TokenDomainValue, securityConfig.CookieSecure, securityConfig.HTTPOnly)
}
