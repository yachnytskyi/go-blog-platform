package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinCookie "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility/cookie"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"
	applicationModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/delivery/http"
	httpGinCommon "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/delivery/http/gin/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location = "internal.user.delivery.http.gin."
	path     = "/"
)

type UserController struct {
	Logger      applicationModel.Logger
	UserUseCase user.UserUseCase
}

func NewUserController(logger applicationModel.Logger, userUseCase user.UserUseCase) UserController {
	return UserController{
		Logger:      logger,
		UserUseCase: userUseCase,
	}
}

func (userController UserController) GetAllUsers(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	paginationQuery := httpGinCommon.ParsePaginationQuery(ginContext)
	fetchedUsers := userController.UserUseCase.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		httpGinCommon.GinNewJSONFailureResponse(ginContext, fetchedUsers.Error, constants.StatusBadRequest)
		return
	}

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.UsersToUsersViewMapper(fetchedUsers.Data)),
	)
}

func (userController UserController) GetCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	currentUser := userController.UserUseCase.GetUserById(ctx, currentUserID)
	if validator.IsError(currentUser.Error) {
		httpGinCommon.GinNewJSONFailureResponse(ginContext, currentUser.Error, constants.StatusBadRequest)
		return
	}

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.UserToUserViewMapper(currentUser.Data)),
	)
}

func (userController UserController) GetUserById(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	userID := ginContext.Param(constants.ItemIdParam)
	fetchedUser := userController.UserUseCase.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		httpGinCommon.GinNewJSONFailureResponse(ginContext, fetchedUser.Error, constants.StatusBadRequest)
		return
	}

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.UserToUserViewMapper(fetchedUser.Data)),
	)
}

func (userController UserController) Register(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userCreateViewData userViewModel.UserCreateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userCreateViewData)
	if validator.IsError(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, userController.Logger, location+"Register", shouldBindJSON)
		return
	}

	userCreateData := userViewModel.UserCreateViewToUserCreateMapper(userCreateViewData)
	createdUser := userController.UserUseCase.Register(ctx, userCreateData)
	if validator.IsError(createdUser.Error) {
		httpGinCommon.GinNewJSONFailureResponse(ginContext, createdUser.Error, constants.StatusBadRequest)
		return
	}

	ginContext.JSON(
		constants.StatusCreated,
		httpModel.NewJSONSuccessResponse(userViewModel.NewWelcomeMessageView(constants.SendingEmailNotification+createdUser.Data.Email)),
	)
}

func (userController UserController) UpdateCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	var userUpdateViewData userViewModel.UserUpdateView
	shouldBindJSON := ginContext.ShouldBindJSON(&userUpdateViewData)
	if validator.IsError(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, userController.Logger, location+"UpdateCurrentUser", shouldBindJSON)
		return
	}

	userUpdateData := userViewModel.UserUpdateViewToUserUpdateMapper(userUpdateViewData)
	userUpdateData.ID = currentUserID
	updatedUser := userController.UserUseCase.UpdateCurrentUser(ctx, userUpdateData)
	if validator.IsError(updatedUser.Error) {
		httpGinCommon.GinNewJSONFailureResponse(ginContext, updatedUser.Error, constants.StatusBadRequest)
		return
	}

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.UserToUserViewMapper(updatedUser.Data)),
	)
}

func (userController UserController) DeleteCurrentUser(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	deletedUser := userController.UserUseCase.DeleteUserById(ctx, currentUserID)
	if validator.IsError(deletedUser) {
		httpGinCommon.GinNewJSONFailureResponse(ginContext, deletedUser, constants.StatusBadRequest)
		return
	}

	httpGinCookie.CleanCookies(ginContext, path)
	ginContext.JSON(constants.StatusNoContent, nil)
}

func (userController UserController) Login(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userLoginViewData userViewModel.UserLoginView
	shouldBindJSON := ginContext.ShouldBindJSON(&userLoginViewData)
	if validator.IsError(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, userController.Logger, location+"Login", shouldBindJSON)
		return
	}

	userLoginData := userViewModel.UserLoginViewToUserLoginMapper(userLoginViewData)
	userToken := userController.UserUseCase.Login(ctx, userLoginData)
	if validator.IsError(userToken.Error) {
		jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(userToken.Error))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	userTokenView := userViewModel.UserTokenToUserTokenViewMapper(userToken.Data)
	setAccessLoginCookies(ginContext, userTokenView.AccessToken, userTokenView.RefreshToken)

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userTokenView),
	)
}

func (userController UserController) RefreshAccessToken(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	currentUserID := ctx.Value(constants.ID).(string)
	currentUser := userController.UserUseCase.GetUserById(ctx, currentUserID)
	if validator.IsError(currentUser.Error) {
		httpGinCommon.GinNewJSONFailureResponse(ginContext, currentUser.Error, constants.StatusBadRequest)
		return
	}

	userToken := userController.UserUseCase.RefreshAccessToken(ctx, currentUser.Data)
	if validator.IsError(userToken.Error) {
		jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(userToken.Error))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	userTokenView := userViewModel.UserTokenToUserTokenViewMapper(userToken.Data)
	setRefreshTokenCookies(ginContext, userTokenView.AccessToken, userTokenView.RefreshToken)

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userTokenView),
	)
}

func (userController UserController) Logout(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	httpGinCookie.CleanCookies(ginContext, path)

	ginContext.JSON(
		constants.StatusOk,
		httpModel.NewJSONSuccessResponse(userViewModel.NewWelcomeMessageView(constants.LogoutNotificationMessage)),
	)
}

func (userController UserController) ForgottenPassword(controllerContext any) {
	ginContext := controllerContext.(*gin.Context)
	ctx, cancel := context.WithTimeout(ginContext.Request.Context(), constants.DefaultContextTimer)
	defer cancel()

	var userForgottenPasswordView userViewModel.UserForgottenPasswordView
	shouldBindJSON := ginContext.ShouldBindJSON(&userForgottenPasswordView)
	if validator.IsError(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, userController.Logger, location+"ForgottenPassword", shouldBindJSON)
		return
	}

	userForgottenPassword := userViewModel.UserForgottenPasswordViewToUserForgottenPassword(userForgottenPasswordView)
	updatedUserError := userController.UserUseCase.ForgottenPassword(ctx, userForgottenPassword)
	if validator.IsError(updatedUserError) {
		jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(updatedUserError))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

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
	shouldBindJSON := ginContext.ShouldBindJSON(&userResetPasswordView)
	if validator.IsError(shouldBindJSON) {
		httpGinCommon.HandleJSONBindingError(ginContext, userController.Logger, location+"ResetUserPassword", shouldBindJSON)
		return
	}

	resetToken := ginContext.Param(constants.ItemIdParam)
	userResetPasswordView.ResetToken = resetToken
	userResetPassword := userViewModel.UserResetPasswordViewToUserResetPassword(userResetPasswordView)

	resetUserPasswordError := userController.UserUseCase.ResetUserPassword(ctx, userResetPassword)
	if validator.IsError(resetUserPasswordError) {
		jsonResponse := httpModel.NewJSONFailureResponse(httpError.HandleError(resetUserPasswordError))
		ginContext.JSON(constants.StatusBadRequest, jsonResponse)
		return
	}

	httpGinCookie.CleanCookies(ginContext, path)
	ginContext.JSON(
		constants.StatusCreated,
		httpModel.NewJSONSuccessResponse(userViewModel.NewWelcomeMessageView(constants.PasswordResetSuccessNotification)),
	)
}

func setAccessLoginCookies(ginContext *gin.Context, accessToken, refreshToken string) {
	accessTokenConfig := config.GetAccessConfig()
	refreshTokenConfig := config.GetRefreshConfig()
	securityConfig := config.GetSecurityConfig()

	ginContext.SetCookie(
		constants.AccessTokenValue,
		accessToken,
		accessTokenConfig.MaxAge,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)

	ginContext.SetCookie(
		constants.RefreshTokenValue,
		refreshToken,
		refreshTokenConfig.MaxAge,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)

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

func setRefreshTokenCookies(ginContext *gin.Context, accessToken, refreshToken string) {
	accessTokenConfig := config.GetAccessConfig()
	refreshTokenConfig := config.GetRefreshConfig()
	securityConfig := config.GetSecurityConfig()

	ginContext.SetCookie(
		constants.AccessTokenValue,
		accessToken,
		accessTokenConfig.MaxAge,
		path,
		constants.TokenDomainValue,
		securityConfig.CookieSecure,
		securityConfig.HTTPOnly,
	)

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
