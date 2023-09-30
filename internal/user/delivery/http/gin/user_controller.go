package gin

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type UserController struct {
	userUseCase user.UserUseCase
}

func NewUserController(userUseCase user.UserUseCase) UserController {
	return UserController{userUseCase: userUseCase}
}

func (userController UserController) GetAllUsers(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), config.DefaultContextTimer)
	defer cancel()
	page := ginContext.DefaultQuery("page", config.DefaultPage)
	limit := ginContext.DefaultQuery("limit", config.DefaultLimit)
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

func (userController UserController) GetMe(ginContext *gin.Context) {
	currentUserView := ginContext.MustGet("currentUser").(userViewModel.UserView)
	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": currentUserView}})
}

func (userController UserController) GetUserById(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), config.DefaultContextTimer)
	defer cancel()
	userID := ginContext.Param("userID")
	defer cancel()
	fetchedUser, err := userController.userUseCase.GetUserById(ctx, userID)
	if validator.IsErrorNotNil(err) {
		if strings.Contains(err.Error(), "Id exists") {
			ginContext.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userViewModel.UserToUserViewMapper(fetchedUser)}})
}

func (userController UserController) Register(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), config.DefaultContextTimer)
	defer cancel()
	var createdUserViewData userViewModel.UserCreateView
	err := ginContext.ShouldBindJSON(&createdUserViewData)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
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

	welcomeMessage := userViewModel.NewWelcomeMessageView(config.SendingEmailNotification + createdUser.Data.Email)
	jsonResponse := httpModel.NewJsonResponseOnSuccess(welcomeMessage)
	httpModel.SetStatus(&jsonResponse)
	ginContext.JSON(http.StatusCreated, jsonResponse)
}

func (userController UserController) UpdateUserById(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), config.DefaultContextTimer)
	defer cancel()
	currentUserView := ginContext.MustGet("currentUser").(userViewModel.UserView)
	userID := currentUserView.UserID
	var updatedUserViewData userViewModel.UserUpdateView

	err := ginContext.ShouldBindJSON(&updatedUserViewData)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedUserData := userViewModel.UserUpdateViewToUserUpdateMapper(updatedUserViewData)
	updatedUser, updatedUserError := userController.userUseCase.UpdateUserById(ctx, userID, updatedUserData)
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

func (userController UserController) Delete(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), config.DefaultContextTimer)
	defer cancel()
	currentUser := ginContext.MustGet("currentUser")
	userID := currentUser.(userViewModel.UserView).UserID
	err := userController.userUseCase.DeleteUser(ctx, fmt.Sprint(userID))
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}
	ginContext.JSON(http.StatusNoContent, nil)
}

func (userController UserController) Login(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), config.DefaultContextTimer)
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

	config, _ := config.LoadConfig(config.ConfigPath)
	accessToken, createTokenError := httpUtility.CreateToken(config.AccessToken.ExpiredIn, userID, config.AccessToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		jsonResponse := httpError.HandleError(createTokenError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	refreshToken, createTokenError := httpUtility.CreateToken(config.RefreshToken.ExpiredIn, userID, config.RefreshToken.PrivateKey)
	if validator.IsErrorNotNil(createTokenError) {
		jsonResponse := httpError.HandleError(createTokenError)
		httpModel.SetStatus(&jsonResponse)
		ginContext.JSON(http.StatusBadRequest, jsonResponse)
		return
	}
	httpGinUtility.LoginSetCookies(ginContext, accessToken, config.AccessToken.MaxAge*60, refreshToken, config.AccessToken.MaxAge*60)
	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userController UserController) RefreshAccessToken(ginContext *gin.Context) {
	message := "could not refresh access token"
	currentUserView := ginContext.MustGet("currentUser").(userViewModel.UserView)

	if validator.IsValueNil(currentUserView) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	cookie, err := ginContext.Cookie("refresh_token")

	if validator.IsErrorNotNil(err) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(config.ConfigPath)
	userID, err := httpUtility.ValidateToken(cookie, config.RefreshToken.PublicKey)
	if validator.IsErrorNotNil(err) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if validator.IsErrorNotNil(err) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user is belonged to this token no longer exists "})
	}

	accessToken, err := httpUtility.CreateToken(config.AccessToken.ExpiredIn, userID, config.AccessToken.PrivateKey)

	if validator.IsErrorNotNil(err) {
		ginContext.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	httpGinUtility.RefreshAccessTokenSetCookies(ginContext, accessToken, config.AccessToken.MaxAge*60)
	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userController UserController) ForgottenPassword(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), config.DefaultContextTimer)
	defer cancel()
	var userViewEmail userViewModel.UserForgottenPasswordView

	err := ginContext.ShouldBindJSON(&userViewEmail)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := config.SendingEmailWithIntstructionsNotifications
	fetchedUser, err := userController.userUseCase.GetUserByEmail(ctx, userViewEmail.Email)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Generate verification code.
	resetToken := randstr.String(20)
	passwordResetToken := common.Encode(resetToken)
	passwordResetAt := time.Now().Add(time.Minute * 15)

	// Update the user.
	err = userController.userUseCase.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Email, "passwordResetToken", passwordResetToken, "passwordResetAt", passwordResetAt)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": err.Error()})
		return
	}
	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (userController UserController) ResetUserPassword(ginContext *gin.Context) {
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), config.DefaultContextTimer)
	defer cancel()
	resetToken := ginContext.Params.ByName("resetToken")
	var userResetPasswordView userViewModel.UserResetPasswordView

	err := ginContext.ShouldBindJSON(&userResetPasswordView)
	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	passwordResetToken := common.Encode(resetToken)

	// Update the user.
	err = userController.userUseCase.ResetUserPassword(ctx, "passwordResetToken", passwordResetToken, "passwordResetAt", "password", userResetPasswordView.Password)

	if validator.IsErrorNotNil(err) {
		ginContext.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	httpGinUtility.ResetUserPasswordSetCookies(ginContext)
	ginContext.JSON(http.StatusOK, gin.H{"status": "success", "message": "Congratulations! Your password was updated successfully! Please sign in again."})

}

func (userController UserController) Logout(ginContext *gin.Context) {
	httpGinUtility.LogoutSetCookies(ginContext)
	ginContext.JSON(http.StatusOK, gin.H{"status": "success"})
}
