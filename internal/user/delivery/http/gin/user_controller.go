package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"

	httpGinCookie "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility/cookie"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	httpGinCommon "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"

	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "internal.user.delivery.http.gin."
	path     = "/"
)

// UserController handles HTTP requests related to user operations.
type UserController struct {
	userUseCase user.UserUseCase
}

// NewUserController creates a new instance of UserController.
func NewUserController(userUseCase user.UserUseCase) UserController {
	return UserController{userUseCase: userUseCase}
}

// GetAllUsers handles an HTTP request to retrieve a paginated list of users.
// It performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Parses the pagination query parameters from the Gin context.
// 3. Calls the user use case to fetch the list of users based on the pagination query.
// 4. Checks for errors during user fetching and responds with a JSON error message if any.
// 5. Maps the fetched user data to a JSON response and sets the HTTP status code to OK.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) GetAllUsers(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Parse pagination query and fetch the users.
	paginationQuery := httpGinCommon.ParsePaginationQuery(ginContext)
	fetchedUsers := userController.userUseCase.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		// Respond with a JSON error message if fetching users fails.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, fetchedUsers.Error, constants.StatusBadRequest)
		return
	}

	// Map the fetched user data to a JSON response and set the status.
	// Return the JSON response with a success status code.
	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.UsersToUsersViewMapper(fetchedUsers.Data)),
	)
}

// GetCurrentUser handles an HTTP request to retrieve the current user's information.
// It performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Retrieves the current user's ID from the context.
// 3. Calls the user use case to fetch the current user's data using the retrieved ID.
// 4. Checks for errors during user fetching and responds with a JSON error message if any.
// 5. Maps the retrieved user data to a JSON response and sets the HTTP status code to OK.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) GetCurrentUser(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Get the current user's ID from the context.
	currentUserID := ctx.Value(constants.ID).(string)

	// Fetch the current user using the user use case.
	currentUser := userController.userUseCase.GetUserById(ctx, currentUserID)
	if validator.IsError(currentUser.Error) {
		// Respond with a JSON error message if fetching the current user fails.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, currentUser.Error, constants.StatusBadRequest)
		return
	}

	// Map the retrieved user data to a JSON response and set the status.
	// Return the JSON response with a success status code.
	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.UserToUserViewMapper(currentUser.Data)),
	)
}

// GetUserById handles an HTTP request to retrieve a user by their ID.
// It performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Retrieves the user ID from the request parameters.
// 3. Calls the user use case to fetch the user's data using the retrieved ID.
// 4. Checks for errors during user fetching and responds with a JSON error message if any.
// 5. Maps the retrieved user data to a JSON response and sets the HTTP status code to OK.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) GetUserById(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Get the user ID from the request parameters.
	userID := ginContext.Param(constants.ItemIdParam)

	// Fetch the user using the user use case.
	fetchedUser := userController.userUseCase.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		// Respond with a JSON error message if fetching the user fails.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, fetchedUser.Error, constants.StatusBadRequest)
		return
	}

	// Map the retrieved user data to a JSON response and set the status.
	// Return the JSON response with a success status code.
	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.UserToUserViewMapper(fetchedUser.Data)),
	)
}

// Register handles an HTTP request to register a new user.
// It performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Binds the incoming JSON data to a struct representing user registration details.
// 3. Responds with a JSON error message if JSON binding fails.
// 4. Maps the user registration view model to a domain model and attempts to register the user.
// 5. Responds with a JSON error message if user registration fails.
// 6. Returns a welcome message in the JSON response upon successful registration.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) Register(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Bind the incoming JSON data to a struct.
	var userCreateViewData userViewModel.UserCreateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userCreateViewData)
	if validator.IsError(shouldBindJSON) {
		// Respond with a JSON error message if JSON binding fails.
		httpGinCommon.HandleJSONBindingError(ginContext, location+"Register", shouldBindJSON)
		return
	}

	// Map the view model to a domain model and attempt user registration.
	userCreateData := userViewModel.UserCreateViewToUserCreateMapper(userCreateViewData)
	createdUser := userController.userUseCase.Register(ctx, userCreateData)
	if validator.IsError(createdUser.Error) {
		// Respond with a JSON error message if user registration fails.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, createdUser.Error, constants.StatusBadRequest)
		return
	}

	// Registration was successful. Return a JSON response with a success status code.
	ginContext.JSON(
		constants.StatusCreated,
		httpModel.NewJSONSuccessResponse(userViewModel.NewWelcomeMessageView(constants.SendingEmailNotification+createdUser.Data.Email)),
	)
}

// UpdateCurrentUser handles an HTTP request to update a user's information based on the provided JSON data.
// It performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Retrieves the current user's ID from the context.
// 3. Binds the incoming JSON data to a struct representing user update details.
// 4. Responds with a JSON error message if JSON binding fails.
// 5. Maps the user update view model to a domain model and updates the user.
// 6. Responds with a JSON error message if user update fails.
// 7. Returns a JSON response with updated user information upon successful update.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) UpdateCurrentUser(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Get the current user's ID from the context.
	currentUserID := ctx.Value(constants.ID).(string)

	// Bind the incoming JSON data to a struct.
	var userUpdateViewData userViewModel.UserUpdateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userUpdateViewData)
	if validator.IsError(shouldBindJSON) {
		// Respond with a JSON error message if JSON binding fails.
		httpGinCommon.HandleJSONBindingError(ginContext, location+"UpdateCurrentUser", shouldBindJSON)
		return
	}

	// Map the view model to a domain model and update the user.
	userUpdateData := userViewModel.UserUpdateViewToUserUpdateMapper(userUpdateViewData)
	userUpdateData.ID = currentUserID
	updatedUser := userController.userUseCase.UpdateCurrentUser(ctx, userUpdateData)
	if validator.IsError(updatedUser.Error) {
		// Respond with a JSON error message if user update fails.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, updatedUser.Error, constants.StatusBadRequest)
		return
	}

	// Update was successful. Return the JSON response with updated user information.
	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.UserToUserViewMapper(updatedUser.Data)),
	)
}

// DeleteCurrentUser handles an HTTP request to delete the current user based on the provided JSON token.
// It performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Retrieves the current user's ID from the context.
// 3. Deletes the user using the user use case.
// 4. Responds with a JSON error message if user deletion fails.
// 5. Cleans cookies to ensure the user is logged out after successful deletion.
// 6. Returns a JSON response with a successful status code (StatusNoContent) indicating successful deletion.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) DeleteCurrentUser(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Extract the current user ID from the context.
	currentUserID := ctx.Value(constants.ID).(string)

	// Delete the user using the user use case.
	deletedUser := userController.userUseCase.DeleteUserById(ctx, currentUserID)
	if validator.IsError(deletedUser) {
		// Respond with a JSON error message if user deletion fails.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, deletedUser, constants.StatusBadRequest)
		return
	}

	// Deletion was successful. Clean cookies to ensure the user is logged out.
	// This prevents any lingering session information after the user has been deleted.
	httpGinCookie.CleanCookies(ginContext, path)

	// Return a JSON response with a successful status code (StatusNoContent).
	ginContext.JSON(constants.StatusNoContent, nil)
}

// Login handles the user authentication process via HTTP POST request.
// It performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Binds the incoming JSON data to a struct representing user login details.
// 3. Handles JSON binding errors and responds with an appropriate error JSON.
// 4. Maps the user login view data to a domain model and authenticates the user using the user use case.
// 5. Responds with a JSON error message if user authentication fails.
// 6. Maps the user token data to a view model and sets cookies for the authenticated user.
// 7. Responds with a success JSON containing the user's access and refresh tokens upon successful authentication.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) Login(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Bind the incoming JSON data to a struct.
	var userLoginViewData userViewModel.UserLoginView
	shouldBindJSON := ginContext.ShouldBindJSON(&userLoginViewData)
	if validator.IsError(shouldBindJSON) {
		// Respond with a JSON error message if JSON binding fails.
		httpGinCommon.HandleJSONBindingError(ginContext, location+"Login", shouldBindJSON)
		return
	}

	// Map the user login view data to a domain model and perform user authentication.
	userLoginData := userViewModel.UserLoginViewToUserLoginMapper(userLoginViewData)
	userToken := userController.userUseCase.Login(ctx, userLoginData)
	if validator.IsError(userToken.Error) {
		// Respond with a JSON error message if user authentication fails.
		jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(userToken.Error))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	// Map user token data to the corresponding user token view data model.
	userTokenView := userViewModel.UserTokenToUserTokenViewMapper(userToken.Data)

	// Set cookies for the authenticated user.
	setLoginCookies(ginContext, userTokenView.AccessToken, userTokenView.RefreshToken)

	// Respond with a success JSON containing the user's access and refresh tokens.
	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userTokenView),
	)
}

// RefreshAccessToken handles the refreshing of the user's access token using the provided refresh token.
// It performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Retrieves the current user's ID from the context and fetches the current user using the user use case.
// 3. Responds with a JSON error message if fetching the current user fails.
// 4. Refreshes the access token using the provided refresh token.
// 5. Responds with a JSON error message if refreshing the access token fails.
// 6. Maps the user token data to a view model and sets cookies with the updated access and refresh tokens.
// 7. Responds with a success JSON containing the user's updated access and refresh tokens upon successful refresh.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) RefreshAccessToken(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Get the current user's ID from the context.
	currentUserID := ctx.Value(constants.ID).(string)

	// Fetch the current user using the user use case.
	currentUser := userController.userUseCase.GetUserById(ctx, currentUserID)
	if validator.IsError(currentUser.Error) {
		// Respond with a JSON error message if fetching the current user fails.
		httpGinCommon.GinNewJSONFailureResponse(ginContext, currentUser.Error, constants.StatusBadRequest)
		return
	}

	// Refresh the access token using the provided refresh token.
	userToken := userController.userUseCase.RefreshAccessToken(ctx, currentUser.Data)
	if validator.IsError(userToken.Error) {
		// Respond with a JSON error message if refreshing the access token fails.
		jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(userToken.Error))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	// Map user token data to the corresponding user token view data model.
	userTokenView := userViewModel.UserTokenToUserTokenViewMapper(userToken.Data)

	// Set cookies with the updated access and refresh tokens.
	setRefreshTokenCookies(ginContext, userTokenView.AccessToken, userTokenView.RefreshToken)

	// Respond with a success JSON containing the user's updated access and refresh tokens.
	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userTokenView),
	)
}

// Logout handles the user logout process.
// It performs the following steps:
// 1. Extracts the Gin context to clean cookies for the logged-out user session.
// 2. Sets cookies to clean, ensuring the user is logged out.
// 3. Responds with a success JSON message indicating successful logout.
//
// Parameters:
// - controllerContext (any): The context object passed from Gin containing HTTP request details.
func (userController UserController) Logout(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)

	// Clean cookies to ensure the user is logged out.
	httpGinCookie.CleanCookies(ginContext, path)

	// Respond with a success JSON message indicating logout.
	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.NewWelcomeMessageView(constants.LogoutNotificationMessage)),
	)
}

// ForgottenPassword handles the request for password recovery.
// It expects a JSON payload containing the user's email and performs the following steps:
// 1. Extracts the Gin context and creates a context with a timeout.
// 2. Binds the incoming JSON data to a struct representing the forgotten password details.
// 3. Handles JSON binding errors if they occur.
// 4. Maps the view model to a domain model for processing.
// 5. Attempts to process the forgotten password request using the user use case.
// 6. Handles errors from the use case and responds with an appropriate error JSON.
// 7. Sends a success response with a welcome message including instructions.
func (userController UserController) ForgottenPassword(controllerContext any) {
	// Extract the Gin context and create a context with a timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Bind the incoming JSON data to a struct representing the forgotten password details.
	var userForgottenPasswordView userViewModel.UserForgottenPasswordView
	shouldBindJSON := ginContext.ShouldBindJSON(&userForgottenPasswordView)
	if validator.IsError(shouldBindJSON) {
		// Handle JSON binding errors.
		httpGinCommon.HandleJSONBindingError(ginContext, location+"ForgottenPassword", shouldBindJSON)
		return
	}

	// Map the view model to a domain model for processing.
	userForgottenPassword := userViewModel.UserForgottenPasswordViewToUserForgottenPassword(userForgottenPasswordView)

	// Attempt to process the forgotten password request using the user use case.
	updatedUserError := userController.userUseCase.ForgottenPassword(ctx, userForgottenPassword)
	if validator.IsError(updatedUserError) {
		// Handle errors from the use case.
		jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(updatedUserError))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	// Send a success response with a welcome message including instructions.
	ginContext.JSON(
		constants.StatusCreated,
		httpModel.NewJSONSuccessResponse(userViewModel.NewWelcomeMessageView(constants.SendingEmailWithInstructionsNotification+" "+userForgottenPassword.Email)),
	)
}

func (userController UserController) ResetUserPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userResetPasswordView userViewModel.UserResetPasswordView

	// Bind the incoming JSON data to a struct.
	shouldBindJSON := ginContext.ShouldBindJSON(&userResetPasswordView)
	if validator.IsError(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, location+"ResetUserPassword", shouldBindJSON)
		return
	}

	// Add the reset token to the struct.
	resetToken := ginContext.Param(constants.ItemIdParam)
	userResetPasswordView.ResetToken = resetToken
	userResetPassword := userViewModel.UserResetPasswordViewToUserResetPassword(userResetPasswordView)

	// Update the user.
	resetUserPasswordError := userController.userUseCase.ResetUserPassword(ctx, userResetPassword)
	if validator.IsError(resetUserPasswordError) {
		// Handle errors from the use case.
		jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(resetUserPasswordError))
		ginContext.JSON(
			constants.StatusBadRequest,
			jsonResponse,
		)
		return
	}

	// Send a success response with a welcome message including instructions.
	httpGinCookie.CleanCookies(ginContext, path)
	ginContext.JSON(
		constants.StatusCreated,
		httpModel.NewJSONSuccessResponse(userViewModel.NewWelcomeMessageView(constants.PasswordResetSuccessNotification)),
	)
}

// setLoginCookies sets cookies in the response for the given access and refresh tokens.
// It performs the following steps:
// 1. Retrieves application configuration for cookie settings.
// 2. Sets the access token cookie with the provided value and configuration.
// 3. Sets the refresh token cookie with the provided value and configuration.
// 4. Sets a "Logged In" flag cookie to indicate the user is authenticated.
func setLoginCookies(ginContext *gin.Context, accessToken, refreshToken string) {
	// Retrieve application configuration for cookie settings.
	accessTokenConfig := config.GetAccessConfig()
	refreshTokenConfig := config.GetRefreshConfig()
	securityConfig := config.GetSecurityConfig()

	// Set the access token cookie with the provided value and configuration.
	ginContext.SetCookie(
		constants.AccessTokenValue,
		accessToken,
		accessTokenConfig.MaxAge,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)

	// Set the refresh token cookie with the provided value and configuration.
	ginContext.SetCookie(
		constants.RefreshTokenValue,
		refreshToken,
		refreshTokenConfig.MaxAge,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)

	// Set the "Logged In" flag cookie to indicate the user is authenticated.
	ginContext.SetCookie(
		constants.LoggedInValue,
		constants.True,
		accessTokenConfig.MaxAge,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)
}

// setRefreshTokenCookies sets cookies in the response for the given access and refresh tokens.
// It performs the following steps:
// 1. Retrieves application configuration for cookie settings.
// 2. Sets the access token cookie with the provided value and configuration.
// 3. Sets the refresh token cookie with the provided value and configuration, if the refresh token is not empty.
// 4. Sets a "Logged In" flag cookie to indicate the user is still authenticated.
func setRefreshTokenCookies(ginContext *gin.Context, accessToken, refreshToken string) {
	// Retrieve application configuration for cookie settings.
	accessTokenConfig := config.GetAccessConfig()
	refreshTokenConfig := config.GetRefreshConfig()
	securityConfig := config.GetSecurityConfig()

	// Set the access token cookie with the provided value and configuration.
	ginContext.SetCookie(
		constants.AccessTokenValue,
		accessToken,
		accessTokenConfig.MaxAge,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)

	// Set the refresh token cookie with the provided value and configuration, if the refresh token is not empty.
	if len(refreshToken) > 0 {
		ginContext.SetCookie(
			constants.RefreshTokenValue,
			refreshToken,
			refreshTokenConfig.MaxAge,
			path,
			constants.TokenDomainValue,
			securityConfig.CookieSecure,
			securityConfig.HTTPOnly,
		)
	}

	// Set the "Logged In" flag cookie to indicate the user is still authenticated.
	ginContext.SetCookie(
		constants.LoggedInValue,
		constants.True,
		accessTokenConfig.MaxAge,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)
}
