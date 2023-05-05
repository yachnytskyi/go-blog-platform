package v1

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	"github.com/yachnytskyi/golang-mongo-grpc/models"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"
)

type UserHandler struct {
	userService user.Service
	template    *template.Template
}

func NewUserHandler(userService user.Service, template *template.Template) UserHandler {
	return UserHandler{userService: userService, template: template}
}

func (userHandler *UserHandler) Register(ctx *gin.Context) {
	var user *models.UserCreate = new(models.UserCreate)
	context := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	createdUser, err := userHandler.userService.Register(context, user)

	if err != nil {
		if strings.Contains(err.Error(), "email already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "error", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": err.Error()})
		return
	}

	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("could not load the config", err)
	}

	// Generate verification code.
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update the user in Database.
	userHandler.userService.UpdateNewRegisteredUserById(context, createdUser.UserID, "verificationCode", verificationCode)

	firstName := createdUser.Name
	firstName = utils.UserFirstName(firstName)

	// Send an email.
	emailData := utils.EmailData{
		URL:       config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(createdUser, &emailData, "verificationCode.html")

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "Error in sending an email"})
		return
	}

	message := "We sent an email with a verification code to " + user.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (userHandler *UserHandler) Login(ctx *gin.Context) {
	var credentials *models.UserSignIn
	context := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	user, err := userHandler.userService.Login(context, credentials)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	config, _ := config.LoadConfig(".")

	// Generate tokens.
	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.UserID, config.AccessTokenPrivateKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	refreshToken, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.UserID, config.RefreshTokenPrivateKey)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	utils.LoginSetCookies(ctx, accessToken, config.AccessTokenMaxAge*60, refreshToken, config.RefreshTokenMaxAge*60)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userHandler *UserHandler) ForgottenPassword(ctx *gin.Context) {
	var userEmail *models.ForgottenPasswordInput
	context := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&userEmail); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	message := "We sent you an email with needed instructions"

	fetchedUser, err := userHandler.userService.GetUserByEmail(context, userEmail.Email)

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
	passwordResetToken := utils.Encode(resetToken)
	passwordResetAt := time.Now().Add(time.Minute * 15)

	// Update the user.
	err = userHandler.userService.UpdatePasswordResetTokenUserByEmail(context, fetchedUser.Email, "passwordResetToken", passwordResetToken, "passwordResetAt", passwordResetAt)

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": err.Error()})
		return
	}

	firstName := fetchedUser.Name
	firstName = utils.UserFirstName(firstName)

	// Send an email.
	emailData := utils.EmailData{
		URL:       config.Origin + "/reset-password/" + resetToken,
		FirstName: firstName,
		Subject:   "Your password reset token (it is valid for 15 minutes)",
	}

	err = utils.SendEmail(fetchedUser, &emailData, "resetPassword.html")

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "success", "message": "Error in sending an email"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": message})

}

func (userHandler *UserHandler) ResetUserPassword(ctx *gin.Context) {
	resetToken := ctx.Params.ByName("resetToken")
	var credentials *models.ResetUserPasswordInput
	context := ctx.Request.Context()

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if credentials.Password != credentials.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	passwordResetToken := utils.Encode(resetToken)

	// Update the user.
	err := userHandler.userService.ResetUserPassword(context, "passwordResetToken", passwordResetToken, "passwordResetAt", "password", credentials.Password)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	utils.ResetUserPasswordSetCookies(ctx)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "Congratulations! Your password was updated successfully! Please sign in again."})

}

func (userHandler *UserHandler) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

	currentUser := ctx.MustGet("currentUser").(*models.User)

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

	userID, err := utils.ValidateToken(cookie, config.RefreshTokenPublicKey)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user is belonged to this token no longer exists "})
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, userID, config.AccessTokenPrivateKey)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	utils.RefreshAccessTokenSetCookies(ctx, accessToken, config.AccessTokenMaxAge*60)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userHandler *UserHandler) UpdateUserById(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.User)
	userID := currentUser.UserID
	context := ctx.Request.Context()

	var updatedUserData *models.UserUpdate = new(models.UserUpdate)

	if err := ctx.ShouldBindJSON(&updatedUserData); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedUser, err := userHandler.userService.UpdateUserById(context, userID, updatedUserData)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "message": "You successfully updated your settings!", "data": gin.H{"user": updatedUser}})
}

func (userHandler *UserHandler) Logout(ctx *gin.Context) {
	utils.LogoutSetCookies(ctx)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (userHandler *UserHandler) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.User)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.UserToUserViewMapper(currentUser)}})
}

func (userHandler *UserHandler) Delete(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser")
	userID := currentUser.(*models.User).UserID
	context := ctx.Request.Context()

	err := userHandler.userService.DeleteUserById(context, fmt.Sprint(userID))

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
	}

	ctx.JSON(http.StatusNoContent, nil)
}
