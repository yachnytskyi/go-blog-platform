package gin

import (
	"context"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"

	httpGinCookie "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility/cookie"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	httpGinCommon "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"

	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "internal.user.delivery.http.gin."
)

type UserController struct {
	userUseCase user.UserUseCase
}

func NewUserController(userUseCase user.UserUseCase) UserController {
	return UserController{userUseCase: userUseCase}
}

// GetAllUsers is a controller method that handles an HTTP request to retrieve a list of users.
// It retrieves user data based on the provided pagination parameters and returns the JSON response.
func (userController UserController) GetAllUsers(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Parse pagination query and fetch the users.
	paginationQuery := httpGinCommon.ParsePaginationQuery(ginContext)
	fetchedUsers := userController.userUseCase.GetAllUsers(ctx, paginationQuery)
	if validator.IsErrorNotNil(fetchedUsers.Error) {
		httpGinCommon.GinNewJSONResponseOnFailure(ginContext, fetchedUsers.Error, constants.StatusBadRequest)
		return
	}

	// Map the fetched user data to a JSON response and set the status.
	// Return the JSON response with a successful status code.
	jsonResponse := httpModel.NewJSONResponseOnSuccess(userViewModel.UsersToUsersViewMapper(fetchedUsers.Data))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

// GetCurrentUser is a controller method that handles an HTTP request to retrieve the current user.
// It retrieves the current user information from a middleware and returns the JSON response.
func (userController UserController) GetCurrentUser(controllerContext any) {
	// Extract the current user from Gin context.
	ginContext := controllerContext.(*gin.Context)
	currentUser := ginContext.MustGet(constants.UserContext).(userViewModel.UserView)
	jsonResponse := httpModel.NewJSONResponseOnSuccess(currentUser)
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

// GetUserById is a controller method that handles an HTTP request to retrieve a user by their ID.
// It retrieves user data and returns the JSON response.
func (userController UserController) GetUserById(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Get the user ID from the request parameters.
	// Fetch the user using the user use case.
	userID := ginContext.Param(constants.UserIDContext)
	fetchedUser := userController.userUseCase.GetUserById(ctx, userID)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		httpGinCommon.GinNewJSONResponseOnFailure(ginContext, fetchedUser.Error, constants.StatusBadRequest)
		return
	}

	// Map the retrieved user data to a JSON response.
	// Return the JSON response with a successful status code.
	jsonResponse := httpModel.NewJSONResponseOnSuccess(userViewModel.UserToUserViewMapper(fetchedUser.Data))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

// Register is a controller method for handling an HTTP request to register a new user.
// It expects a controller context and performs the following steps:
// 1. Binds the incoming JSON data to a struct representing user registration details.
// 2. If the JSON binding fails, it responds with an error message.
// 3. Converts the user registration view model to a user model and attempts user registration.
// 4. If the registration fails, it responds with an error message.
// 5. If registration is successful, it returns a welcome message in the JSON response.
func (userController UserController) Register(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Bind the incoming JSON data to a struct.
	var userCreateViewData userViewModel.UserCreateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userCreateViewData)
	if validator.IsErrorNotNil(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, location+"Register", shouldBindJSON)
		return
	}

	// Convert the view model to a domain model and register the user.
	userCreateData := userViewModel.UserCreateViewToUserCreateMapper(userCreateViewData)
	createdUser := userController.userUseCase.Register(ctx, userCreateData)
	if validator.IsErrorNotNil(createdUser.Error) {
		httpGinCommon.GinNewJSONResponseOnFailure(ginContext, createdUser.Error, constants.StatusBadRequest)
		return
	}

	// Registration was successful. Return the JSON response with a successful status code.
	welcomeMessage := userViewModel.NewWelcomeMessageView(constants.SendingEmailNotification + createdUser.Data.Email)
	jsonResponse := httpModel.NewJSONResponseOnSuccess(welcomeMessage)
	ginContext.JSON(constants.StatusCreated, jsonResponse)
}

// UpdateUserById updates a user's information based on the provided JSON data
// and returns the JSON response.
func (userController UserController) UpdateCurrentUser(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Get the current user's ID from the Gin context.
	// Bind the incoming JSON data to a struct.
	currentUserID := ginContext.MustGet(constants.UserIDContext).(string)
	var userUpdateViewData userViewModel.UserUpdateView

	// Bind the incoming JSON data to a struct.
	shouldBindJSON := ginContext.ShouldBindJSON(&userUpdateViewData)
	if validator.IsErrorNotNil(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, location+"UpdateCurrentUser", shouldBindJSON)
		return
	}

	// Convert the view model to a domain model and update the user.
	userUpdateData := userViewModel.UserUpdateViewToUserUpdateMapper(userUpdateViewData)
	userUpdateData.UserID = currentUserID
	updatedUser := userController.userUseCase.UpdateCurrentUser(ctx, userUpdateData)
	if validator.IsErrorNotNil(updatedUser.Error) {
		httpGinCommon.GinNewJSONResponseOnFailure(ginContext, updatedUser.Error, constants.StatusBadRequest)
		return
	}

	// Update was successful. Return the JSON response with a successful status code.
	jsonResponse := httpModel.NewJSONResponseOnSuccess(userViewModel.UserToUserViewMapper(updatedUser.Data))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

// Delete deletes a user based on the provided JSON token.
func (userController UserController) DeleteCurrentUser(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Extract the current user ID from the gin context.
	// Delete the user using the user use case.
	currentUserID := ginContext.MustGet(constants.UserIDContext).(string)
	deletedUser := userController.userUseCase.DeleteUserById(ctx, currentUserID)
	if validator.IsErrorNotNil(deletedUser) {
		httpGinCommon.GinNewJSONResponseOnFailure(ginContext, deletedUser, constants.StatusBadRequest)
		return
	}

	// Deletion was successful. Clean cookies to ensure the user is logged out, and return a JSON response
	// with a successful status code (StatusNoContent). Cleaning cookies is necessary to prevent any
	// lingering session information after the user has been deleted.
	httpGinCookie.CleanCookies(ginContext)
	ginContext.JSON(constants.StatusNoContent, nil)
}

// Login handles the user authentication process.
// It expects a JSON payload containing user login information.
// Upon successful authentication, it sets cookies and responds with a success JSON.
// In case of errors, it handles them and responds with an appropriate error JSON.
func (userController UserController) Login(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Bind the incoming JSON data to a struct.
	var userLoginViewData userViewModel.UserLoginView
	shouldBindJSON := ginContext.ShouldBindJSON(&userLoginViewData)
	if validator.IsErrorNotNil(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, location+"Login", shouldBindJSON)
		return
	}

	// Map the user login view data to the internal user login data model and perform user authentication using the user use case.
	userLoginData := userViewModel.UserLoginViewToUserLoginMapper(userLoginViewData)
	userLogin := userController.userUseCase.Login(ctx, userLoginData)
	if validator.IsErrorNotNil(userLogin.Error) {
		jsonResponse := httpModel.NewJSONResponseOnFailure(httpError.HandleError(userLogin.Error))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	// Set cookies for the authenticated user.
	setLoginCookies(ginContext, userLogin.Data.AccessToken, userLogin.Data.RefreshToken)

	// Respond with a success JSON containing the user's access token.
	jsonResponse := httpModel.NewJSONResponseOnSuccess(userViewModel.TokenStringToTokenViewMapper(userLogin.Data.AccessToken))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

func (userController UserController) RefreshAccessToken(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	cookie, cookieError := ginContext.Cookie(constants.RefreshTokenValue)
	if validator.IsErrorNotNil(cookieError) {
		jsonResponse := httpModel.NewJSONResponseOnFailure(httpError.HandleError(cookieError))
		ginContext.JSON(constants.StatusUnauthorized, jsonResponse)
		return
	}

	accessTokenConfig := config.AppConfig.AccessToken
	refreshTokenConfig := config.AppConfig.RefreshToken
	userID, validateTokenError := domainUtility.ValidateJWTToken(cookie, refreshTokenConfig.PublicKey)
	if validator.IsErrorNotNil(validateTokenError) {
		jsonResponse := httpModel.NewJSONResponseOnFailure(httpError.HandleError(validateTokenError))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	accessToken, createTokenError := domainUtility.GenerateJWTToken(ctx, accessTokenConfig.ExpiredIn, userID, accessTokenConfig.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		jsonResponse := httpModel.NewJSONResponseOnFailure(httpError.HandleError(createTokenError))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	// Generate a new refresh token (optional, based on your requirements).
	newRefreshToken, newRefreshTokenError := domainUtility.GenerateJWTToken(ctx, refreshTokenConfig.ExpiredIn, userID, refreshTokenConfig.PrivateKey)
	if validator.IsErrorNotNil(newRefreshTokenError) {
		jsonResponse := httpModel.NewJSONResponseOnFailure(httpError.HandleError(newRefreshTokenError))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}
	setRefreshTokenCookies(ginContext, accessToken, newRefreshToken)
	jsonResponse := httpModel.NewJSONResponseOnSuccess(userViewModel.TokenStringToTokenViewMapper(accessToken))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

func (userController UserController) Logout(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	httpGinCookie.CleanCookies(ginContext)
	ginContext.JSON(constants.StatusOk, gin.H{"status": "success"})
}

func (userController UserController) ForgottenPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	var userViewEmail userViewModel.UserForgottenPasswordView

	// Bind the incoming JSON data to a struct.
	shouldBindJSON := ginContext.ShouldBindJSON(&userViewEmail)
	if validator.IsErrorNotNil(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, location+"ForgottenPassword", shouldBindJSON)
		return
	}

	fetchedUser := userController.userUseCase.GetUserByEmail(ctx, userViewEmail.Email)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		ginContext.JSON(constants.StatusBadGateway, gin.H{"status": "error", "message": fetchedUser.Error})
		return
	}

	// Generate verification code.
	resetToken := randstr.String(20)
	passwordResetToken := commonUtility.Encode(resetToken)
	passwordResetAt := time.Now().Add(time.Minute * 15)

	// Update the user.
	updatedUser := userController.userUseCase.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Data.Email, "passwordResetToken", passwordResetToken, "passwordResetAt", passwordResetAt)
	if validator.IsErrorNotNil(updatedUser) {
		ginContext.JSON(constants.StatusBadGateway, gin.H{"status": "success", "message": updatedUser.Error()})
		return
	}
	ginContext.JSON(constants.StatusOk, gin.H{"status": "success", "message": constants.SendingEmailWithInstructionsNotification})
}

func (userController UserController) ResetUserPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	resetToken := ginContext.Params.ByName("resetToken")
	var userResetPasswordView userViewModel.UserResetPasswordView

	// Bind the incoming JSON data to a struct.
	shouldBindJSON := ginContext.ShouldBindJSON(&userResetPasswordView)
	if validator.IsErrorNotNil(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, location+"ResetUserPassword", shouldBindJSON)
		return
	}

	passwordResetToken := commonUtility.Encode(resetToken)

	// Update the user.
	err := userController.userUseCase.ResetUserPassword(ctx, "passwordResetToken", passwordResetToken, "passwordResetAt", "password", userResetPasswordView.Password)

	if validator.IsErrorNotNil(err) {
		ginContext.JSON(constants.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	httpGinCookie.CleanCookies(ginContext)
	ginContext.JSON(constants.StatusOk, gin.H{"status": "success", "message": "Congratulations! Your password was updated successfully! Please sign in again."})

}

// setLoginCookies sets cookies in the response for the given access and refresh tokens.
// It is typically used during user login to store authentication-related information.
func setLoginCookies(ginContext *gin.Context, accessToken, refreshToken string) {
	// Retrieve application configuration for cookie settings.
	accessTokenConfig := config.AppConfig.AccessToken
	refreshTokenConfig := config.AppConfig.RefreshToken
	securityConfig := config.AppConfig.Security

	// Set the access token cookie with the provided value and configuration.
	ginContext.SetCookie(constants.AccessTokenValue, accessToken, accessTokenConfig.MaxAge, "/", constants.TokenDomainValue,
		securityConfig.CookieSecure, true)

	// Set the refresh token cookie with the provided value and configuration.
	ginContext.SetCookie(constants.RefreshTokenValue, refreshToken, refreshTokenConfig.MaxAge, "/", constants.TokenDomainValue,
		securityConfig.CookieSecure, true)

	// Set the "Logged In" flag cookie to indicate the user is authenticated.
	ginContext.SetCookie(constants.LoggedInValue, constants.True, accessTokenConfig.MaxAge, "/", constants.TokenDomainValue,
		securityConfig.CookieSecure, false)
}

// setRefreshTokenCookies sets cookies in the response for the given access and refresh tokens.
// It is typically used when refreshing the access token.
func setRefreshTokenCookies(ginContext *gin.Context, accessToken, refreshToken string) {
	// Retrieve application configuration for cookie settings.
	accessTokenConfig := config.AppConfig.AccessToken
	refreshTokenConfig := config.AppConfig.RefreshToken
	securityConfig := config.AppConfig.Security

	// Set the access token cookie with the provided value and configuration.
	ginContext.SetCookie(constants.AccessTokenValue, accessToken, accessTokenConfig.MaxAge, "/", constants.TokenDomainValue,
		securityConfig.CookieSecure, true)

	// Set the refresh token cookie with the provided value and configuration (optional).
	if validator.IsStringNotEmpty(refreshToken) {
		ginContext.SetCookie(constants.RefreshTokenValue, refreshToken, refreshTokenConfig.MaxAge, "/", constants.TokenDomainValue,
			securityConfig.CookieSecure, true)
	}

	// Set the "Logged In" flag cookie to indicate the user is still logged in.
	ginContext.SetCookie(constants.LoggedInValue, constants.True, accessTokenConfig.MaxAge, "/", constants.TokenDomainValue,
		securityConfig.CookieSecure, false)
}
