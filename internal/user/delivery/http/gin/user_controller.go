package gin

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"

	httpGinCookie "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility/cookie"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	httpGinCommon "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"

	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type UserController struct {
	userUseCase user.UserUseCase
}

func NewUserController(userUseCase user.UserUseCase) UserController {
	return UserController{userUseCase: userUseCase}
}

// GetAllUsers is a controller method for handling an HTTP request to retrieve a list of users.
// It expects a controller context and retrieves user data based on the provided pagination parameters.
//
// Possible Errors:
// - http.StatusBadRequest: If there is an issue with the request parameters or the user data retrieval.
// - Other errors depending on userUseCase.GetAllUsers.

func (userController UserController) GetAllUsers(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	paginationQuery := httpGinCommon.ParsePaginationQuery(ginContext)
	fetchedUsers := userController.userUseCase.GetAllUsers(ctx, paginationQuery)
	if validator.IsErrorNotNil(fetchedUsers.Error) {
		httpGinCommon.GinNewJsonResponseOnFailure(ginContext, fetchedUsers.Error, http.StatusBadRequest)
		return
	}

	// Map the fetched user data to a JSON response and set the status.
	// Send the JSON response with a successful status code.
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UsersToUsersViewMapper(fetchedUsers.Data))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusOK, jsonResponse)
}

func (userController UserController) GetCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(ginContext.MustGet(constants.UserContext).(userViewModel.UserView))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusOK, jsonResponse)
}

// GetUserById is a controller method for handling an HTTP request to retrieve a user by their ID.
// It expects a controller context and the user's ID to fetch the corresponding user from the database.
func (userController UserController) GetUserById(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	// Get the user ID from the request parameters.
	// Fetch the user using the user use case.
	userID := ginContext.Param(constants.UserIDContext)
	fetchedUser := userController.userUseCase.GetUserById(ctx, userID)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		httpGinCommon.GinNewJsonResponseOnFailure(ginContext, fetchedUser.Error, http.StatusBadRequest)
		return
	}

	// Map the retrieved user data to a JSON response.
	// Return the JSON response with a successful status code.
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UserToUserViewMapper(fetchedUser.Data))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusOK, jsonResponse)
}

func (userController UserController) Register(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	var createdUserViewData userViewModel.UserCreateView
	shouldBindJSON := ginContext.ShouldBindJSON(&createdUserViewData)
	if validator.IsErrorNotNil(shouldBindJSON) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": shouldBindJSON.Error()})
		return
	}
	userCreate := userViewModel.UserCreateViewToUserCreateMapper(createdUserViewData)
	createdUser := userController.userUseCase.Register(ctx, userCreate)
	if validator.IsErrorNotNil(createdUser.Error) {
		createdUserError := httpError.HandleError(createdUser.Error)
		jsonResponse := httpModel.NewJsonResponseOnFailure(createdUserError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	welcomeMessage := userViewModel.NewWelcomeMessageView(constants.SendingEmailNotification + createdUser.Data.Email)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(welcomeMessage)
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusCreated, jsonResponse)
}

func (userController UserController) UpdateUserById(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	currentUserID := ginContext.MustGet(constants.UserIDContext).(string)

	var updatedUserViewData userViewModel.UserUpdateView
	err := ginContext.ShouldBindJSON(&updatedUserViewData)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedUserData := userViewModel.UserUpdateViewToUserUpdateMapper(updatedUserViewData)
	updatedUserData.UserID = currentUserID
	updatedUser, updatedUserError := userController.userUseCase.UpdateUserById(ctx, updatedUserData)
	if validator.IsErrorNotNil(updatedUserError) {
		updatedUserError := httpError.HandleError(updatedUserError)
		jsonResponse := httpModel.NewJsonResponseOnFailure(updatedUserError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UserToUserViewMapper(updatedUser))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusCreated, jsonResponse)
}

func (userController UserController) Delete(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	currentUserID := ginContext.MustGet(constants.UserIDContext).(string)
	err := userController.userUseCase.DeleteUser(ctx, currentUserID)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}
	ginContext.JSON(http.StatusNoContent, nil)
}

func (userController UserController) Login(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	var userLoginViewData userViewModel.UserLoginView
	err := ginContext.ShouldBindJSON(&userLoginViewData)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userLoginData := userViewModel.UserLoginViewToUserLoginMapper(userLoginViewData)
	userID, loginError := userController.userUseCase.Login(ctx, userLoginData)
	if validator.IsErrorNotNil(loginError) {
		loginError := httpError.HandleError(loginError)
		jsonResponse := httpModel.NewJsonResponseOnFailure(loginError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	applicationConfig := config.AppConfig
	accessToken, createTokenError := httpUtility.CreateToken(applicationConfig.AccessToken.ExpiredIn, userID, applicationConfig.AccessToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		createTokenError := httpError.HandleError(createTokenError)
		jsonResponse := httpModel.NewJsonResponseOnFailure(createTokenError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	refreshToken, createTokenError := httpUtility.CreateToken(applicationConfig.RefreshToken.ExpiredIn, userID, applicationConfig.RefreshToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		createTokenError := httpError.HandleError(createTokenError)
		jsonResponse := httpModel.NewJsonResponseOnFailure(createTokenError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	ginContext.SetCookie(constants.AccessTokenValue, accessToken, applicationConfig.AccessToken.MaxAge, "/", constants.TokenDomainValue, false, true)
	ginContext.SetCookie(constants.RefreshTokenValue, refreshToken, applicationConfig.RefreshToken.MaxAge, "/", constants.TokenDomainValue, false, true)
	ginContext.SetCookie(constants.LoggedInValue, "true", applicationConfig.AccessToken.MaxAge, "/", constants.TokenDomainValue, false, false)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.TokenStringToTokenViewMapper(accessToken))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusCreated, jsonResponse)
}

func (userController UserController) RefreshAccessToken(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	message := "could not refresh access token"
	currentUser := ginContext.MustGet(constants.UserContext).(userViewModel.UserView)

	if validator.IsValueNil(currentUser) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}
	cookie, err := ginContext.Cookie(constants.RefreshTokenValue)

	if validator.IsErrorNotNil(err) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	applicationConfig := config.AppConfig
	userID, validateTokenError := httpUtility.ValidateToken(cookie, applicationConfig.RefreshToken.PublicKey)
	if validator.IsErrorNotNil(validateTokenError) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": validateTokenError.Error()})
		return
	}

	if validator.IsErrorNotNil(err) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user is belonged to this token no longer exists "})
	}

	accessToken, createTokenError := httpUtility.CreateToken(applicationConfig.AccessToken.ExpiredIn, userID, applicationConfig.AccessToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}
	ginContext.SetCookie(constants.AccessTokenValue, accessToken, applicationConfig.AccessToken.MaxAge, "/", constants.TokenDomainValue, false, true)
	ginContext.SetCookie(constants.LoggedInValue, "true", applicationConfig.AccessToken.MaxAge, "/", constants.TokenDomainValue, false, false)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.TokenStringToTokenViewMapper(accessToken))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusCreated, jsonResponse)
}

func (userController UserController) ForgottenPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	var userViewEmail userViewModel.UserForgottenPasswordView

	err := ginContext.ShouldBindJSON(&userViewEmail)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fetchedUser, err := userController.userUseCase.GetUserByEmail(ctx, userViewEmail.Email)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Generate verification code.
	resetToken := randstr.String(20)
	passwordResetToken := commonUtility.Encode(resetToken)
	passwordResetAt := time.Now().Add(time.Minute * 15)

	// Update the user.
	err = userController.userUseCase.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Email, "passwordResetToken", passwordResetToken, "passwordResetAt", passwordResetAt)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "message": constants.SendingEmailWithInstructionsNotification})
}

func (userController UserController) ResetUserPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	resetToken := ginContext.Params.ByName("resetToken")
	var userResetPasswordView userViewModel.UserResetPasswordView

	err := ginContext.ShouldBindJSON(&userResetPasswordView)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	passwordResetToken := commonUtility.Encode(resetToken)

	// Update the user.
	err = userController.userUseCase.ResetUserPassword(ctx, "passwordResetToken", passwordResetToken, "passwordResetAt", "password", userResetPasswordView.Password)

	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	httpGinCookie.CleanCookies(ginContext)
	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "message": "Congratulations! Your password was updated successfully! Please sign in again."})

}

func (userController UserController) Logout(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	httpGinCookie.CleanCookies(ginContext)
	ginContext.JSON(http.StatusOK, gin.H{"status": "success"})
}
