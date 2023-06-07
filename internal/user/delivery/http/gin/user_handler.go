package gin

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	httpGinUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/gin/utility"
	httpUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/http/utility"
	userViewModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/delivery/model"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
)

type UserHandler struct {
	userUseCase user.UseCase
	template    *template.Template
}

func NewUserHandler(userUseCase user.UseCase, template *template.Template) UserHandler {
	return UserHandler{userUseCase: userUseCase, template: template}
}

func (userHandler *UserHandler) GetAllUsers(ctx *gin.Context) {
	page := ctx.DefaultQuery("page", "1")
	limit := ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	fetchedUsers, err := userHandler.userUseCase.GetAllUsers(ctx.Request.Context(), intPage, intLimit)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, userViewModel.UsersToUsersViewMapper(fetchedUsers))
}

func (userHandler *UserHandler) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*userModel.User)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userViewModel.UserToUserViewMapper(currentUser)}})
}

func (userHandler *UserHandler) GetUserById(ctx *gin.Context) {
	userID := ctx.Param("userID")

	fetchedUser, err := userHandler.userUseCase.GetUserById(ctx.Request.Context(), userID)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": userViewModel.UserToUserViewMapper(fetchedUser)}})
}

func (userHandler *UserHandler) Register(ctx *gin.Context) {
	var createdUserViewData *userViewModel.UserCreateView = new(userViewModel.UserCreateView)

	if err := ctx.ShouldBindJSON(&createdUserViewData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := createdUserViewData.UserCreateViewValidator(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	createdUserData := userViewModel.UserCreateViewToUserCreateMapper((*userViewModel.UserCreateView)(createdUserViewData))
	createdUser, err := userHandler.userUseCase.Register(ctx.Request.Context(), &createdUserData)

	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load the config", err)
	}

	// Generate verification code.
	code := randstr.String(20)

	verificationCode := utility.Encode(code)

	// Update the user in Database.
	userHandler.userUseCase.UpdateNewRegisteredUserById(ctx.Request.Context(), createdUser.UserID, "verificationCode", verificationCode)

	firstName := createdUser.Name
	firstName = httpUtility.UserFirstName(firstName)

	// Send an email.
	emailData := httpUtility.EmailData{
		URL:       config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = httpUtility.SendEmail(createdUser, &emailData, "verificationCode.html")

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "Error in sending an email"})
		return
	}

	message := "We sent an email with a verification code to " + createdUserData.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (userHandler *UserHandler) UpdateUserById(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*userModel.User)
	userID := currentUser.UserID

	var updatedUserViewData *userViewModel.UserUpdateView = new(userViewModel.UserUpdateView)

	if err := ctx.ShouldBindJSON(&updatedUserViewData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := updatedUserViewData.UserUpdateViewValidator(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedUserData := userViewModel.UserUpdateViewToUserUpdateMapper(updatedUserViewData)
	updatedUser, err := userHandler.userUseCase.UpdateUserById(ctx.Request.Context(), userID, &updatedUserData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully updated your settings!", "data": gin.H{"user": userViewModel.UserToUserViewMapper(updatedUser)}})
}

func (userHandler *UserHandler) Delete(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	userID := currentUser.(*userModel.User).UserID

	err := userHandler.userUseCase.DeleteUserById(ctx.Request.Context(), fmt.Sprint(userID))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	ctx.JSON(http.StatusNoContent, nil)
}

func (userHandler *UserHandler) Login(ctx *gin.Context) {
	var userLoginViewData *userViewModel.UserLoginView

	if err := ctx.ShouldBindJSON(&userLoginViewData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err := userLoginViewData.UserSignInViewValidator(); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	userLoginData := userViewModel.UserLoginViewToUserLoginMapper(userLoginViewData)
	userID, err := userHandler.userUseCase.Login(ctx.Request.Context(), &userLoginData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	config, _ := config.LoadConfig(".")

	// Generate tokens.
	accessToken, err := httpUtility.CreateToken(config.AccessTokenExpiresIn, userID, config.AccessTokenPrivateKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := httpUtility.CreateToken(config.RefreshTokenExpiresIn, userID, config.RefreshTokenPrivateKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	httpGinUtility.LoginSetCookies(ctx, accessToken, config.AccessTokenMaxAge*60, refreshToken, config.RefreshTokenMaxAge*60)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userHandler *UserHandler) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	currentUser := ctx.MustGet("currentUser").(*userModel.User)

	if currentUser == nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	cookie, err := ctx.Cookie("refresh_token")

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": message})
		return
	}

	config, _ := config.LoadConfig(".")

	userID, err := httpUtility.ValidateToken(cookie, config.RefreshTokenPublicKey)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user is belonged to this token no longer exists "})
	}

	accessToken, err := httpUtility.CreateToken(config.AccessTokenExpiresIn, userID, config.AccessTokenPrivateKey)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	httpGinUtility.RefreshAccessTokenSetCookies(ctx, accessToken, config.AccessTokenMaxAge*60)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userHandler *UserHandler) ForgottenPassword(ctx *gin.Context) {
	var userEmail *userModel.UserForgottenPassword

	if err := ctx.ShouldBindJSON(&userEmail); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "We sent you an email with needed instructions"

	fetchedUser, err := userHandler.userUseCase.GetUserByEmail(ctx.Request.Context(), userEmail.Email)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	if !fetchedUser.Verified {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": "Please verify you account"})
		return
	}

	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	// Generate verification code.
	resetToken := randstr.String(20)
	passwordResetToken := utility.Encode(resetToken)
	passwordResetAt := time.Now().Add(time.Minute * 15)

	// Update the user.
	err = userHandler.userUseCase.UpdatePasswordResetTokenUserByEmail(ctx.Request.Context(), fetchedUser.Email, "passwordResetToken", passwordResetToken, "passwordResetAt", passwordResetAt)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": err.Error()})
		return
	}

	firstName := fetchedUser.Name
	firstName = httpUtility.UserFirstName(firstName)

	// Send an email.
	emailData := httpUtility.EmailData{
		URL:       config.Origin + "/reset-password/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (it is valid for 15 minutes)",
	}

	err = httpUtility.SendEmail(fetchedUser, &emailData, "resetPassword.html")

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "Error in sending an email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})
}

func (userHandler *UserHandler) ResetUserPassword(ctx *gin.Context) {
	resetToken := ctx.Params.ByName("resetToken")
	var credentials *userModel.UserResetPassword

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if credentials.Password != credentials.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	passwordResetToken := utility.Encode(resetToken)

	// Update the user.
	err := userHandler.userUseCase.ResetUserPassword(ctx.Request.Context(), "passwordResetToken", passwordResetToken, "passwordResetAt", "password", credentials.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	httpGinUtility.ResetUserPasswordSetCookies(ctx)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Congratulations! Your password was updated successfully! Please sign in again."})

}

func (userHandler *UserHandler) Logout(ctx *gin.Context) {
	httpGinUtility.LogoutSetCookies(ctx)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}
