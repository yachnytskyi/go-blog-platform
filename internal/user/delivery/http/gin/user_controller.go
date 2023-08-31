package gin

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/model"

	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	httpError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/http"

	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
	httpModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/model/http"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

type UserController struct {
	userUseCase user.UseCase
}

func NewUserController(userUseCase user.UseCase) UserController {
	return UserController{userUseCase: userUseCase}
}

func (userController *UserController) GetAllUsers(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fetchedUsers, err := userController.userUseCase.GetAllUsers(ctx.Request.Context(), intPage, intLimit)
	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	jsonResponse := httpModel.NewJsonResponse(userViewModel.UsersToUsersViewMapper(fetchedUsers))
	httpModel.SetStatus(jsonResponse)
	ctx.JSON(http.StatusOK, jsonResponse)

}

func (userController *UserController) GetMe(ctx *gin.Context) {
	currentUserView := ctx.MustGet("currentUser").(*userViewModel.UserView)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": currentUserView}})
}

func (userController *UserController) GetUserById(ctx *gin.Context) {
	userID := ctx.Param("userID")
	fetchedUser, err := userController.userUseCase.GetUserById(ctx.Request.Context(), userID)
	if validator.IsErrorNotNil(err) {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userViewModel.UserToUserViewMapper(fetchedUser)}})
}

func (userController *UserController) Register(ctx *gin.Context) {
	var createdUserViewData *userViewModel.UserCreateView = new(userViewModel.UserCreateView)

	err := ctx.ShouldBindJSON(&createdUserViewData)
	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	createdUserData := userViewModel.UserCreateViewToUserCreateMapper(createdUserViewData)
	createdUserError := userController.userUseCase.Register(ctx.Request.Context(), createdUserData)

	if createdUserError != nil {
		jsonResponse := httpError.HandleError(createdUserError)
		httpModel.SetStatus(jsonResponse)
		ctx.JSON(http.StatusBadRequest, jsonResponse)
		return
	}

	welcomeMessage := &userViewModel.UserWelcomeMessageView{
		Message: config.SendingEmailNotification + createdUserData.Email,
	}
	jsonResponse := httpModel.NewJsonResponse(welcomeMessage)
	httpModel.SetStatus(jsonResponse)
	ctx.JSON(http.StatusCreated, jsonResponse)
}

func (userController *UserController) UpdateUserById(ctx *gin.Context) {
	currentUserView := ctx.MustGet("currentUser").(*userViewModel.UserView)
	userID := currentUserView.UserID
	var updatedUserViewData *userViewModel.UserUpdateView = new(userViewModel.UserUpdateView)

	err := ctx.ShouldBindJSON(&updatedUserViewData)
	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedUserData := userViewModel.UserUpdateViewToUserUpdateMapper(updatedUserViewData)
	updatedUser, updatedUserError := userController.userUseCase.UpdateUserById(ctx.Request.Context(), userID, &updatedUserData)

	if updatedUserError != nil {
		ctx.JSON(http.StatusBadRequest, httpError.HandleError(updatedUserError))
	}

	ctx.JSON(http.StatusOK, userViewModel.UserToUserViewMapper(updatedUser))
}

func (userController *UserController) Delete(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	userID := currentUser.(*userViewModel.UserView).UserID
	err := userController.userUseCase.DeleteUserById(ctx.Request.Context(), fmt.Sprint(userID))

	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (userController *UserController) Login(ctx *gin.Context) {
	var userLoginViewData *userViewModel.UserLoginView

	err := ctx.ShouldBindJSON(&userLoginViewData)
	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userLoginData := userViewModel.UserLoginViewToUserLoginMapper(userLoginViewData)
	userID, err := userController.userUseCase.Login(ctx.Request.Context(), &userLoginData)

	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	config, _ := config.LoadConfig(".")

	// Generate tokens.
	accessToken, err := httpUtility.CreateToken(config.AccessTokenExpiresIn, userID, config.AccessTokenPrivateKey)

	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := httpUtility.CreateToken(config.RefreshTokenExpiresIn, userID, config.RefreshTokenPrivateKey)

	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	httpGinUtility.LoginSetCookies(ctx, accessToken, config.AccessTokenMaxAge*60, refreshToken, config.RefreshTokenMaxAge*60)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userController *UserController) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"
	currentUserView := ctx.MustGet("currentUser").(*userViewModel.UserView)

	if currentUserView == nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	cookie, err := ctx.Cookie("refresh_token")

	if validator.IsErrorNotNil(err) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(".")

	userID, err := httpUtility.ValidateToken(cookie, config.RefreshTokenPublicKey)

	if validator.IsErrorNotNil(err) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if validator.IsErrorNotNil(err) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user is belonged to this token no longer exists "})
	}

	accessToken, err := httpUtility.CreateToken(config.AccessTokenExpiresIn, userID, config.AccessTokenPrivateKey)

	if validator.IsErrorNotNil(err) {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	httpGinUtility.RefreshAccessTokenSetCookies(ctx, accessToken, config.AccessTokenMaxAge*60)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userController *UserController) ForgottenPassword(ctx *gin.Context) {
	var userViewEmail *userViewModel.UserForgottenPasswordView

	err := ctx.ShouldBindJSON(&userViewEmail)
	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := config.SendingEmailWithIntstructionsNotifications
	fetchedUser, err := userController.userUseCase.GetUserByEmail(ctx.Request.Context(), userViewEmail.Email)

	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// Generate verification code.
	resetToken := randstr.String(20)
	passwordResetToken := utility.Encode(resetToken)
	passwordResetAt := time.Now().Add(time.Minute * 15)

	// Update the user.
	err = userController.userUseCase.UpdatePasswordResetTokenUserByEmail(ctx.Request.Context(), fetchedUser.Email, "passwordResetToken", passwordResetToken, "passwordResetAt", passwordResetAt)

	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (userController *UserController) ResetUserPassword(ctx *gin.Context) {
	resetToken := ctx.Params.ByName("resetToken")
	var userResetPasswordView *userViewModel.UserResetPasswordView

	err := ctx.ShouldBindJSON(&userResetPasswordView)
	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	passwordResetToken := utility.Encode(resetToken)

	// Update the user.
	err = userController.userUseCase.ResetUserPassword(ctx.Request.Context(), "passwordResetToken", passwordResetToken, "passwordResetAt", "password", userResetPasswordView.Password)

	if validator.IsErrorNotNil(err) {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	httpGinUtility.ResetUserPasswordSetCookies(ctx)
	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Congratulations! Your password was updated successfully! Please sign in again."})

}

func (userController *UserController) Logout(ctx *gin.Context) {
	httpGinUtility.LogoutSetCookies(ctx)
	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
