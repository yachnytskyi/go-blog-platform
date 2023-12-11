package cookie

import (
	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
)

// CleanCookies is a helper function for Gin controllers that clears specific cookies.
func CleanCookies(ctx *gin.Context) {
	// Clear the access token cookie by setting its value to an empty string and
	// configuring it with the logout max age, path, and domain.
	ctx.SetCookie(constants.AccessTokenValue, constants.EmptyString, constants.LogoutMaxAgeValue, "/", constants.TokenDomainValue, false, true)

	// Clear the refresh token cookie in a similar manner.
	ctx.SetCookie("refresh_token", constants.EmptyString, constants.LogoutMaxAgeValue, "/", constants.TokenDomainValue, false, true)

	// Clear the "loggedIn" cookie in the same way.
	ctx.SetCookie(constants.LoggedInValue, constants.EmptyString, constants.LogoutMaxAgeValue, "/", constants.TokenDomainValue, false, true)
}
