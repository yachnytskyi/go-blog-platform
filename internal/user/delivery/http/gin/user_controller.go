package gin

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	utility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility/cookie"
	view "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	domain "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/usecase"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	delivery "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "internal.user.delivery.http.gin."
	path     = "/"
)

type UserController struct {
	Config      *config.ApplicationConfig
	Logger      interfaces.Logger
	UserUseCase domain.UserUseCase
}

func NewUserController(config *config.ApplicationConfig, logger interfaces.Logger, userUseCase domain.UserUseCase) UserController {
	return UserController{
		Config:      config,
		Logger:      logger,
		UserUseCase: userUseCase,
	}
}

func (userController UserController) GetAllUsers(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	paginationQuery := common.ParsePaginationQuery(ginContext)
	fetchedUsers := userController.UserUseCase.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(fetchedUsers.Error)))
		return
	}

	ginContext.JSON(
		http.StatusOK, model.NewJSONResponseOnSuccess(view.UsersToUsersViewMapper(fetchedUsers.Data)))
}

func (userController UserController) GetCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	currentUser := userController.UserUseCase.GetUserById(ctx, currentUserID)
	if validator.IsError(currentUser.Error) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(currentUser.Error)))
		return
	}

	ginContext.JSON(http.StatusOK, model.NewJSONResponseOnSuccess(view.UserToUserViewMapper(currentUser.Data)))
}

func (userController UserController) GetUserById(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	userID := ginContext.Param(constants.ItemIdParam)
	fetchedUser := userController.UserUseCase.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(fetchedUser.Error)))
		return
	}

	ginContext.JSON(http.StatusOK, model.NewJSONResponseOnSuccess(view.UserToUserViewMapper(fetchedUser.Data)))
}

func (userController UserController) Register(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userCreateViewData view.UserCreateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userCreateViewData)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"Register", shouldBindJSON)
		return
	}

	userCreateData := view.UserCreateViewToUserCreateMapper(userCreateViewData)
	createdUser := userController.UserUseCase.Register(ctx, userCreateData)
	if validator.IsError(createdUser.Error) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(createdUser.Error)))
		return
	}

	ginContext.JSON(
		http.StatusCreated,
		model.NewJSONResponseOnSuccess(view.NewWelcomeMessageView(fmt.Sprintf(constants.SendingEmailNotification, createdUser.Data.Email))),
	)
}

func (userController UserController) UpdateCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	var userUpdateViewData view.UserUpdateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userUpdateViewData)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"UpdateCurrentUser", shouldBindJSON)
		return
	}

	userUpdateData := view.UserUpdateViewToUserUpdateMapper(userUpdateViewData)
	userUpdateData.ID = currentUserID
	updatedUser := userController.UserUseCase.UpdateCurrentUser(ctx, userUpdateData)
	if validator.IsError(updatedUser.Error) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(updatedUser.Error)))
		return
	}

	ginContext.JSON(http.StatusOK, model.NewJSONResponseOnSuccess(view.UserToUserViewMapper(updatedUser.Data)))
}

func (userController UserController) DeleteCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	deletedUserError := userController.UserUseCase.DeleteUserById(ctx, currentUserID)
	if validator.IsError(deletedUserError) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(deletedUserError)))
		return
	}

	utility.CleanCookies(ginContext, userController.Config, path)
	ginContext.JSON(http.StatusNoContent, nil)
}

func (userController UserController) Login(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userLoginViewData view.UserLoginView
	shouldBindJSON := ginContext.ShouldBindJSON(&userLoginViewData)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"Login", shouldBindJSON)
		return
	}

	userLoginData := view.UserLoginViewToUserLoginMapper(userLoginViewData)
	userToken := userController.UserUseCase.Login(ctx, userLoginData)
	if validator.IsError(userToken.Error) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(userToken.Error)))
		return
	}

	userTokenView := view.UserTokenToUserTokenViewMapper(userToken.Data)
	setAccessLoginCookies(ginContext, userController.Config, userTokenView.AccessToken, userTokenView.RefreshToken)
	ginContext.JSON(http.StatusOK, model.NewJSONResponseOnSuccess(userTokenView))
}

func (userController UserController) RefreshAccessToken(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	currentUser := userController.UserUseCase.GetUserById(ctx, currentUserID)
	if validator.IsError(currentUser.Error) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(currentUser.Error)))
		return
	}

	userToken := userController.UserUseCase.RefreshAccessToken(ctx, currentUser.Data)
	if validator.IsError(userToken.Error) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(userToken.Error)))
		return
	}

	userTokenView := view.UserTokenToUserTokenViewMapper(userToken.Data)
	setRefreshTokenCookies(ginContext, userController.Config, userTokenView.AccessToken, userTokenView.RefreshToken)
	ginContext.JSON(http.StatusOK, model.NewJSONResponseOnSuccess(userTokenView))
}

func (userController UserController) Logout(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	utility.CleanCookies(ginContext, userController.Config, path)
	ginContext.JSON(http.StatusOK, model.NewJSONResponseOnSuccess(view.NewWelcomeMessageView(constants.LogoutNotificationMessage)))
}

func (userController UserController) ForgottenPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userForgottenPasswordView view.UserForgottenPasswordView
	shouldBindJSON := ginContext.ShouldBindJSON(&userForgottenPasswordView)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"ForgottenPassword", shouldBindJSON)
		return
	}

	userForgottenPassword := view.UserForgottenPasswordViewToUserForgottenPassword(userForgottenPasswordView)
	updatedUserError := userController.UserUseCase.ForgottenPassword(ctx, userForgottenPassword)
	if validator.IsError(updatedUserError) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(updatedUserError)))
		return
	}

	ginContext.JSON(
		http.StatusCreated,
		model.NewJSONResponseOnSuccess(view.NewWelcomeMessageView(constants.SendingEmailWithInstructionsNotification+" "+userForgottenPassword.Email)))
}

func (userController UserController) ResetUserPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userResetPasswordView view.UserResetPasswordView
	shouldBindJSON := ginContext.ShouldBindJSON(&userResetPasswordView)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"ResetUserPassword", shouldBindJSON)
		return
	}

	resetToken := ginContext.Param(constants.ItemIdParam)
	userResetPasswordView.ResetToken = resetToken
	userResetPassword := view.UserResetPasswordViewToUserResetPassword(userResetPasswordView)
	resetUserPasswordError := userController.UserUseCase.ResetUserPassword(ctx, userResetPassword)
	if validator.IsError(resetUserPasswordError) {
		ginContext.JSON(http.StatusBadRequest, model.NewJSONResponseOnFailure(delivery.HandleError(resetUserPasswordError)))
		return
	}

	utility.CleanCookies(ginContext, userController.Config, path)
	ginContext.JSON(http.StatusCreated, model.NewJSONResponseOnSuccess(view.NewWelcomeMessageView(constants.PasswordResetSuccessNotification)))
}

func setAccessLoginCookies(ginContext *gin.Context, config *config.ApplicationConfig, accessToken, refreshToken string) {
	ginContext.SetCookie(
		constants.AccessTokenValue,
		accessToken,
		config.AccessToken.MaxAge,
		path,
		constants.TokenDomainValue,
		config.Security.CookieSecure,
		config.Security.HTTPOnly,
	)

	ginContext.SetCookie(
		constants.RefreshTokenValue,
		refreshToken,
		config.RefreshToken.MaxAge,
		path,
		constants.TokenDomainValue,
		config.Security.CookieSecure,
		config.Security.HTTPOnly,
	)

	ginContext.SetCookie(
		constants.LoggedInValue,
		constants.True,
		config.AccessToken.MaxAge,
		path,
		constants.TokenDomainValue,
		config.Security.CookieSecure,
		config.Security.HTTPOnly,
	)
}

func setRefreshTokenCookies(ginContext *gin.Context, config *config.ApplicationConfig, accessToken, refreshToken string) {
	ginContext.SetCookie(
		constants.AccessTokenValue,
		accessToken,
		config.AccessToken.MaxAge,
		path,
		constants.TokenDomainValue,
		config.Security.CookieSecure,
		config.Security.HTTPOnly,
	)

	if refreshToken != "" {
		ginContext.SetCookie(
			constants.RefreshTokenValue,
			refreshToken,
			config.RefreshToken.MaxAge,
			path,
			constants.TokenDomainValue,
			config.Security.CookieSecure,
			config.Security.HTTPOnly,
		)
	}

	ginContext.SetCookie(
		constants.LoggedInValue,
		constants.True,
		config.AccessToken.MaxAge,
		path,
		constants.TokenDomainValue,
		config.Security.CookieSecure,
		config.Security.HTTPOnly,
	)
}
