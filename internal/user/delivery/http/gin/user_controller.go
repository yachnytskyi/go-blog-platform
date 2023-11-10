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
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"

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

	paginationQuery := httpGinCommon.ParsePaginationQuery(ginContext)
	fetchedUsers := userController.userUseCase.GetAllUsers(ctx, paginationQuery)
	if validator.IsErrorNotNil(fetchedUsers.Error) {
		httpGinCommon.GinNewJsonResponseOnFailure(ginContext, fetchedUsers.Error, constants.StatusBadRequest)
		return
	}

	// Map the fetched user data to a JSON response and set the status.
	// Return the JSON response with a successful status code.
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UsersToUsersViewMapper(fetchedUsers.Data))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

// GetCurrentUser is a controller method that handles an HTTP request to retrieve the current user.
// It extracts the current user information from a middleware and returns the JSON response.
func (userController UserController) GetCurrentUser(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	currentUser := ginContext.MustGet(constants.UserContext).(userViewModel.UserView)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(currentUser)
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
		httpGinCommon.GinNewJsonResponseOnFailure(ginContext, fetchedUser.Error, constants.StatusBadRequest)
		return
	}

	// Map the retrieved user data to a JSON response.
	// Return the JSON response with a successful status code.
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UserToUserViewMapper(fetchedUser.Data))
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
	var createdUserViewData userViewModel.UserCreateView
	shouldBindJSON := ginContext.ShouldBindJSON(&createdUserViewData)
	if validator.IsErrorNotNil(shouldBindJSON) {
		internalError := httpError.NewHttpInternalErrorView(location+"Register.ShouldBindJSON", shouldBindJSON.Error())
		logging.Logger(internalError)
		httpGinCommon.GinNewJsonResponseOnFailure(ginContext, internalError, constants.StatusBadRequest)
		return
	}

	// Convert the view model to a domain model and register the user.
	userCreate := userViewModel.UserCreateViewToUserCreateMapper(createdUserViewData)
	createdUser := userController.userUseCase.Register(ctx, userCreate)
	if validator.IsErrorNotNil(createdUser.Error) {
		httpGinCommon.GinNewJsonResponseOnFailure(ginContext, createdUser.Error, constants.StatusBadRequest)
		return
	}

	// Registration was successful. Return the JSON response with a successful status code.
	welcomeMessage := userViewModel.NewWelcomeMessageView(constants.SendingEmailNotification + createdUser.Data.Email)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(welcomeMessage)
	ginContext.JSON(constants.StatusCreated, jsonResponse)
}

// UpdateUserById updates a user's information based on the provided JSON data.
func (userController UserController) UpdateUserById(controllerContext any) {
	// Extract the Gin context and create a context with timeout.
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Get the current user's ID from the Gin context.
	// Bind the incoming JSON data to a struct.
	currentUserID := ginContext.MustGet(constants.UserIDContext).(string)
	var updatedUserViewData userViewModel.UserUpdateView
	shouldBindJSON := ginContext.ShouldBindJSON(&updatedUserViewData)
	if validator.IsErrorNotNil(shouldBindJSON) {
		internalError := httpError.NewHttpInternalErrorView(location+"UpdateUserById.ShouldBindJSON", shouldBindJSON.Error())
		logging.Logger(internalError)
		httpGinCommon.GinNewJsonResponseOnFailure(ginContext, internalError, constants.StatusBadRequest)
		return
	}

	// Convert the view model to a domain model and update the user.
	updatedUserData := userViewModel.UserUpdateViewToUserUpdateMapper(updatedUserViewData)
	updatedUserData.UserID = currentUserID
	updatedUser, updateErr := userController.userUseCase.UpdateUserById(ctx, updatedUserData)
	if validator.IsErrorNotNil(updateErr) {
		httpGinCommon.GinNewJsonResponseOnFailure(ginContext, updateErr, constants.StatusBadRequest)
		return
	}

	// Return the JSON response with a successful status code.
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UserToUserViewMapper(updatedUser))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

func (userController UserController) Delete(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	currentUserID := ginContext.MustGet(constants.UserIDContext).(string)
	shouldBindJSON := userController.userUseCase.DeleteUser(ctx, currentUserID)
	if validator.IsErrorNotNil(shouldBindJSON) {
		ginContext.JSON(constants.StatusBadRequest, gin.H{"status": "fail", "message": shouldBindJSON.Error()})
	}
	ginContext.JSON(constants.StatusNoContent, nil)
}

func (userController UserController) Login(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	var userLoginViewData userViewModel.UserLoginView
	shouldBindJSON := ginContext.ShouldBindJSON(&userLoginViewData)
	if validator.IsErrorNotNil(shouldBindJSON) {
		ginContext.JSON(constants.StatusBadRequest, gin.H{"status": "fail", "message": shouldBindJSON.Error()})
		return
	}

	userLoginData := userViewModel.UserLoginViewToUserLoginMapper(userLoginViewData)
	userID, loginError := userController.userUseCase.Login(ctx, userLoginData)
	if validator.IsErrorNotNil(loginError) {
		loginError := httpError.HandleError(loginError)
		jsonResponse := httpModel.NewJsonResponseOnFailure(loginError)
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	applicationConfig := config.AppConfig
	accessToken, createTokenError := domainUtility.GenerateJWTToken(applicationConfig.AccessToken.ExpiredIn, userID, applicationConfig.AccessToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		createTokenError := httpError.HandleError(createTokenError)
		jsonResponse := httpModel.NewJsonResponseOnFailure(createTokenError)
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}
	refreshToken, createTokenError := domainUtility.GenerateJWTToken(applicationConfig.RefreshToken.ExpiredIn, userID, applicationConfig.RefreshToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		createTokenError := httpError.HandleError(createTokenError)
		jsonResponse := httpModel.NewJsonResponseOnFailure(createTokenError)
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}
	ginContext.SetCookie(constants.AccessTokenValue, accessToken, applicationConfig.AccessToken.MaxAge, "/", constants.TokenDomainValue, false, true)
	ginContext.SetCookie(constants.RefreshTokenValue, refreshToken, applicationConfig.RefreshToken.MaxAge, "/", constants.TokenDomainValue, false, true)
	ginContext.SetCookie(constants.LoggedInValue, "true", applicationConfig.AccessToken.MaxAge, "/", constants.TokenDomainValue, false, false)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.TokenStringToTokenViewMapper(accessToken))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

func (userController UserController) RefreshAccessToken(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	message := "could not refresh access token"
	currentUser := ginContext.MustGet(constants.UserContext).(userViewModel.UserView)

	if validator.IsValueNil(currentUser) {
		ginContext.AbortWithStatusJSON(constants.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}
	cookie, cookieError := ginContext.Cookie(constants.RefreshTokenValue)

	if validator.IsErrorNotNil(cookieError) {
		ginContext.AbortWithStatusJSON(constants.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	applicationConfig := config.AppConfig
	userID, validateTokenError := domainUtility.ValidateJWTToken(cookie, applicationConfig.RefreshToken.PublicKey)
	if validator.IsErrorNotNil(validateTokenError) {
		ginContext.AbortWithStatusJSON(constants.StatusForbidden, gin.H{"status": "fail", "message": validateTokenError.Error()})
		return
	}

	if validator.IsErrorNotNil(cookieError) {
		ginContext.AbortWithStatusJSON(constants.StatusForbidden, gin.H{"status": "fail", "message": "the user is belonged to this token no longer exists "})
	}

	accessToken, createTokenError := domainUtility.GenerateJWTToken(applicationConfig.AccessToken.ExpiredIn, userID, applicationConfig.AccessToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		ginContext.AbortWithStatusJSON(constants.StatusForbidden, gin.H{"status": "fail", "message": cookieError.Error()})
		return
	}
	ginContext.SetCookie(constants.AccessTokenValue, accessToken, applicationConfig.AccessToken.MaxAge, "/", constants.TokenDomainValue, false, true)
	ginContext.SetCookie(constants.LoggedInValue, "true", applicationConfig.AccessToken.MaxAge, "/", constants.TokenDomainValue, false, false)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.TokenStringToTokenViewMapper(accessToken))
	ginContext.JSON(constants.StatusOk, jsonResponse)
}

func (userController UserController) ForgottenPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	var userViewEmail userViewModel.UserForgottenPasswordView

	shouldBindJSON := ginContext.ShouldBindJSON(&userViewEmail)
	if validator.IsErrorNotNil(shouldBindJSON) {
		ginContext.JSON(constants.StatusBadRequest, gin.H{"status": "fail", "message": shouldBindJSON.Error()})
		return
	}

	fetchedUser := userController.userUseCase.GetUserByEmail(ctx, userViewEmail.Email)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		ginContext.JSON(constants.StatusBadGateway, gin.H{"status": "error", "message": shouldBindJSON.Error()})
		return
	}

	// Generate verification code.
	resetToken := randstr.String(20)
	passwordResetToken := commonUtility.Encode(resetToken)
	passwordResetAt := time.Now().Add(time.Minute * 15)

	// Update the user.
	shouldBindJSON = userController.userUseCase.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Data.Email, "passwordResetToken", passwordResetToken, "passwordResetAt", passwordResetAt)
	if validator.IsErrorNotNil(shouldBindJSON) {
		ginContext.JSON(constants.StatusBadGateway, gin.H{"status": "success", "message": shouldBindJSON.Error()})
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

	shouldBindJSON := ginContext.ShouldBindJSON(&userResetPasswordView)
	if validator.IsErrorNotNil(shouldBindJSON) {
		ginContext.JSON(constants.StatusBadRequest, gin.H{"status": "fail", "message": shouldBindJSON.Error()})
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

func (userController UserController) Logout(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	httpGinCookie.CleanCookies(ginContext)
	ginContext.JSON(constants.StatusOk, gin.H{"status": "success"})
}
