package cookie

import (
	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

// CleanCookies is a helper function for Gin controllers that clears specific cookies.
// Parameters:
//   - ctx: The Gin context which provides access to the request and response objects.
//   - path: The path attribute for the cookies to be cleared. This specifies the URL path
//     for which the cookies are valid. The cookies will be cleared for this path.
func CleanCookies(ctx *gin.Context, path string) {
	// Retrieve application configuration for cookie settings.
	securityConfig := config.GetSecurityConfig()

	// Clear the access token cookie by setting its value to an empty string and
	// configuring it with the logout max age, path, and domain.
	// The security settings from the configuration are applied to ensure proper handling.
	ctx.SetCookie(
		constants.AccessTokenValue,  // Name of the cookie
		"",                          // Value of the cookie (empty to clear)
		constants.LogoutMaxAgeValue, // Max age of the cookie
		path,                        // Path for which the cookie is valid
		constants.TokenDomainValue,  // Domain for which the cookie is valid
		securityConfig.CookieSecure, // Secure flag from configuration
		securityConfig.HTTPOnly,     // HTTPOnly flag from configuration
	)

	// Clear the refresh token cookie in a similar manner.
	ctx.SetCookie(
		constants.RefreshTokenValue,
		"",
		constants.LogoutMaxAgeValue,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)

	// Clear the "loggedIn" cookie in the same way.
	ctx.SetCookie(
		constants.LoggedInValue,
		"",
		constants.LogoutMaxAgeValue,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)
}
