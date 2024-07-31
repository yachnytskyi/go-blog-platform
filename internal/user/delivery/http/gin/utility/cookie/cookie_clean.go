package cookie

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
)

// CleanCookies is a helper function for Gin controllers that clears specific cookies.
func CleanCookies(ctx *gin.Context, configInstance interfaces.Config, path string) {
	config := configInstance.GetConfig()

	// Clear the access token cookie by setting its value to an empty string and
	// configuring it with the logout max age, path, and domain.
	// The security settings from the configuration are applied to ensure proper handling.
	ctx.SetCookie(
		constants.AccessTokenValue,   // Name of the cookie
		"",                           // Value of the cookie (empty to clear)
		constants.LogoutMaxAgeValue,  // Max age of the cookie
		path,                         // Path for which the cookie is valid
		constants.TokenDomainValue,   // Domain for which the cookie is valid
		config.Security.CookieSecure, // Secure flag from configuration
		config.Security.HTTPOnly,     // HTTPOnly flag from configuration
	)

	// Clear the refresh token cookie in a similar manner.
	ctx.SetCookie(
		constants.RefreshTokenValue,
		"",
		constants.LogoutMaxAgeValue,
		path,
		constants.TokenDomainValue,
		config.Security.CookieSecure,
		config.Security.HTTPOnly,
	)

	// Clear the "loggedIn" cookie in the same way.
	ctx.SetCookie(
		constants.LoggedInValue,
		"",
		constants.LogoutMaxAgeValue,
		path,
		constants.TokenDomainValue,
		config.Security.CookieSecure,
		config.Security.HTTPOnly,
	)
}
