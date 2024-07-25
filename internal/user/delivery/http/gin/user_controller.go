package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinCookie "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility/cookie"
	userView "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "internal.user.delivery.http.gin."
	path     = "/"
)

type UserController struct {
	Config      model.Config
	Logger      model.Logger
	UserUseCase user.UserUseCase
}

func NewUserController(config model.Config, logger model.Logger, userUseCase user.UserUseCase) user.UserController {
	return &UserController{
		Config:      config,
		Logger:      logger,
		UserUseCase: userUseCase,
	}
}

func (userController *UserController) GetAllUsers(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	paginationQuery := common.ParsePaginationQuery(ginContext)
	fetchedUsers := userController.UserUseCase.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(fetchedUsers.Error)))
		return
	}

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONResponseOnSuccess(userView.UsersToUsersViewMapper(fetchedUsers.Data)),
	)
}

func (userController *UserController) GetCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	currentUser := userController.UserUseCase.GetUserById(ctx, currentUserID)
	if validator.IsError(currentUser.Error) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(currentUser.Error)))
		return
	}

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONResponseOnSuccess(userView.UserToUserViewMapper(currentUser.Data)),
	)
}

func (userController *UserController) GetUserById(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	userID := ginContext.Param(constants.ItemIdParam)
	fetchedUser := userController.UserUseCase.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(fetchedUser.Error)))
		return
	}

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONResponseOnSuccess(userView.UserToUserViewMapper(fetchedUser.Data)),
	)
}

func (userController *UserController) Register(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userCreateViewData userView.UserCreateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userCreateViewData)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"Register", shouldBindJSON)
		return
	}

	userCreateData := userView.UserCreateViewToUserCreateMapper(userCreateViewData)
	createdUser := userController.UserUseCase.Register(ctx, userCreateData)
	if validator.IsError(createdUser.Error) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(createdUser.Error)))
		return
	}

	ginContext.JSON(
		constants.StatusCreated,
		httpModel.NewJSONResponseOnSuccess(userView.NewWelcomeMessageView(constants.SendingEmailNotification+createdUser.Data.Email)),
	)
}

func (userController *UserController) UpdateCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	var userUpdateViewData userView.UserUpdateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userUpdateViewData)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"UpdateCurrentUser", shouldBindJSON)
		return
	}

	userUpdateData := userView.UserUpdateViewToUserUpdateMapper(userUpdateViewData)
	userUpdateData.ID = currentUserID
	updatedUser := userController.UserUseCase.UpdateCurrentUser(ctx, userUpdateData)
	if validator.IsError(updatedUser.Error) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(updatedUser.Error)))
		return
	}

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONResponseOnSuccess(userView.UserToUserViewMapper(updatedUser.Data)),
	)
}

func (userController *UserController) DeleteCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	deletedUserError := userController.UserUseCase.DeleteUserById(ctx, currentUserID)
	if validator.IsError(deletedUserError) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(deletedUserError)))
		return
	}

	httpGinCookie.CleanCookies(ginContext, userController.Config, path)
	ginContext.JSON(constants.StatusNoContent, nil)
}

func (userController *UserController) Login(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userLoginViewData userView.UserLoginView
	shouldBindJSON := ginContext.ShouldBindJSON(&userLoginViewData)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"Login", shouldBindJSON)
		return
	}

	userLoginData := userView.UserLoginViewToUserLoginMapper(userLoginViewData)
	userToken := userController.UserUseCase.Login(ctx, userLoginData)
	if validator.IsError(userToken.Error) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(userToken.Error)))
		return
	}

	userTokenView := userView.UserTokenToUserTokenViewMapper(userToken.Data)
	setAccessLoginCookies(ginContext, userController.Config, userTokenView.AccessToken, userTokenView.RefreshToken)

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONResponseOnSuccess(userTokenView),
	)
}

func (userController *UserController) RefreshAccessToken(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	currentUser := userController.UserUseCase.GetUserById(ctx, currentUserID)
	if validator.IsError(currentUser.Error) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(currentUser.Error)))
		return
	}

	userToken := userController.UserUseCase.RefreshAccessToken(ctx, currentUser.Data)
	if validator.IsError(userToken.Error) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(userToken.Error)))
		return
	}

	userTokenView := userView.UserTokenToUserTokenViewMapper(userToken.Data)
	setRefreshTokenCookies(ginContext, userController.Config, userTokenView.AccessToken, userTokenView.RefreshToken)

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONResponseOnSuccess(userTokenView),
	)
}

func (userController *UserController) Logout(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	httpGinCookie.CleanCookies(ginContext, userController.Config, path)

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONResponseOnSuccess(userView.NewWelcomeMessageView(constants.LogoutNotificationMessage)),
	)
}

func (userController *UserController) ForgottenPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userForgottenPasswordView userView.UserForgottenPasswordView
	shouldBindJSON := ginContext.ShouldBindJSON(&userForgottenPasswordView)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"ForgottenPassword", shouldBindJSON)
		return
	}

	userForgottenPassword := userView.UserForgottenPasswordViewToUserForgottenPassword(userForgottenPasswordView)
	updatedUserError := userController.UserUseCase.ForgottenPassword(ctx, userForgottenPassword)
	if validator.IsError(updatedUserError) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(updatedUserError)))
		return
	}

	ginContext.JSON(
		constants.StatusCreated,
		httpModel.NewJSONResponseOnSuccess(userView.NewWelcomeMessageView(constants.SendingEmailWithInstructionsNotification+" "+userForgottenPassword.Email)),
	)
}

func (userController *UserController) ResetUserPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userResetPasswordView userView.UserResetPasswordView
	shouldBindJSON := ginContext.ShouldBindJSON(&userResetPasswordView)
	if validator.IsError(shouldBindJSON) {
		common.HandleJSONBindingError(ginContext, userController.Logger, location+"ResetUserPassword", shouldBindJSON)
		return
	}

	resetToken := ginContext.Param(constants.ItemIdParam)
	userResetPasswordView.ResetToken = resetToken
	userResetPassword := userView.UserResetPasswordViewToUserResetPassword(userResetPasswordView)

	resetUserPasswordError := userController.UserUseCase.ResetUserPassword(ctx, userResetPassword)
	if validator.IsError(resetUserPasswordError) {
		ginContext.JSON(constants.StatusBadRequest, httpModel.NewJSONResponseOnFailure(httpError.HandleError(resetUserPasswordError)))
		return
	}

	httpGinCookie.CleanCookies(ginContext, userController.Config, path)
	ginContext.JSON(
		constants.StatusCreated,
		httpModel.NewJSONResponseOnSuccess(userView.NewWelcomeMessageView(constants.PasswordResetSuccessNotification)),
	)
}

func setAccessLoginCookies(ginContext *gin.Context, configInstance model.Config, accessToken, refreshToken string) {
	config := configInstance.GetConfig()

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

func setRefreshTokenCookies(ginContext *gin.Context, configInstance model.Config, accessToken, refreshToken string) {
	config := configInstance.GetConfig()

	ginContext.SetCookie(
		constants.AccessTokenValue,
		accessToken,
		config.AccessToken.MaxAge,
		path,
		constants.TokenDomainValue,
		config.Security.CookieSecure,
		config.Security.HTTPOnly,
	)

	if len(refreshToken) > 0 {
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
