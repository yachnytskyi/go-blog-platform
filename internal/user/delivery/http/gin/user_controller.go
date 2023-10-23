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
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
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

func (userController UserController) GetAllUsers(controllerContext interface{}) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	page := ginContext.DefaultQuery("page", constants.DefaultPage)
	limit := ginContext.DefaultQuery("limit", constants.DefaultLimit)
	orderBy := ginContext.DefaultQuery("order-by", "")
	convertedPage := commonModel.GetPage(page)
	convertedLimit := commonModel.GetLimit(limit)
	paginationQuery := commonModel.NewPaginationQuery(convertedPage, convertedLimit, orderBy)
	fetchedUsers := userController.userUseCase.GetAllUsers(ctx, paginationQuery)
	if validator.IsErrorNotNil(fetchedUsers.Error) {
		jsonResponse := httpError.HandleError(fetchedUsers.Error)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UsersToUsersViewMapper(fetchedUsers.Data))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusOK, jsonResponse)
}

func (userController UserController) GetCurrentUser(controllerContext interface{}) {
	ginContext := controllerContext.(*gin.Context)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(ginContext.MustGet(constants.UserContext).(userViewModel.UserView))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusOK, jsonResponse)
}

func (userController UserController) GetUserById(controllerContext interface{}) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()
	userID := ginContext.Param(constants.UserIDContext)
	fetchedUser := userController.userUseCase.GetUserById(ctx, userID)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		jsonResponse := httpError.HandleError(fetchedUser.Error)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UserToUserViewMapper(fetchedUser.Data))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusOK, jsonResponse)
}

func (userController UserController) Register(controllerContext interface{}) {
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
		jsonResponse := httpError.HandleError(createdUser.Error)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	welcomeMessage := userViewModel.NewWelcomeMessageView(constants.SendingEmailNotification + createdUser.Data.Email)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(welcomeMessage)
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusCreated, jsonResponse)
}

func (userController UserController) UpdateUserById(controllerContext interface{}) {
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
	updatedUser, updatedUserError := userController.userUseCase.UpdateUserById(ctx, currentUserID, updatedUserData)
	if validator.IsErrorNotNil(updatedUserError) {
		jsonResponse := httpError.HandleError(updatedUserError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	jsonResponse := httpModel.NewJsonResponseOnSuccess(userViewModel.UserToUserViewMapper(updatedUser))
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusCreated, jsonResponse)
}

func (userController UserController) Delete(controllerContext interface{}) {
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

func (userController UserController) Login(controllerContext interface{}) {
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
		jsonResponse := httpError.HandleError(loginError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	applicationConfig := config.AppConfig
	accessToken, createTokenError := httpUtility.CreateToken(applicationConfig.AccessToken.ExpiredIn, userID, applicationConfig.AccessToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		jsonResponse := httpError.HandleError(createTokenError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	refreshToken, createTokenError := httpUtility.CreateToken(applicationConfig.RefreshToken.ExpiredIn, userID, applicationConfig.RefreshToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		jsonResponse := httpError.HandleError(createTokenError)
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

func (userController UserController) RefreshAccessToken(controllerContext interface{}) {
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

func (userController UserController) ForgottenPassword(controllerContext interface{}) {
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

func (userController UserController) ResetUserPassword(controllerContext interface{}) {
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

func (userController UserController) Logout(controllerContext interface{}) {
	ginContext := controllerContext.(*gin.Context)
	httpGinCookie.CleanCookies(ginContext)
	ginContext.JSON(http.StatusOK, gin.H{"status": "success"})
}
