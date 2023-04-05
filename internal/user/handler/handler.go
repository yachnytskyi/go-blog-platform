package handler

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

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

func NewUserHandler(userService user.Service, template *template.Template) user.Handler {
	return &UserHandler{userService: userService, template: template}
}

func (userHandler *UserHandler) Register(ctx *gin.Context) {
	var user *models.UserCreate

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	if user.Password != user.PasswordConfirm {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": "Passwords do not match"})
		return
	}

	context := ctx.Request.Context()

	newUser, err := userHandler.userService.Register(context, user)

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

	// Generate Verification Code.
	code := randstr.String(20)

	verificationCode := utils.Encode(code)

	// Update the user in Database.
	userHandler.userService.UpdateUserById(context, newUser.UserID.Hex(), "verificationCode", verificationCode)

	firstName := newUser.Name

	if strings.Contains(firstName, " ") {
		firstName = strings.Split(firstName, " ")[1]
	}

	// firstName = utils.UserFirstName(firstName)

	// Send an email.
	emailData := utils.EmailData{
		URL:       config.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	err = utils.SendEmail(newUser, &emailData, userHandler.template, "verificationCode.html")

	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": "Error in sending an email"})
		return
	}

	message := "We sent an email with a verification code to " + user.Email
	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "message": message})
}

func (userHandler *UserHandler) Login(ctx *gin.Context) {
	var credentials *models.UserSignIn

	if err := ctx.ShouldBindJSON(&credentials); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	context := ctx.Request.Context()

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

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", refreshToken, config.RefreshTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userHandler *UserHandler) RefreshAccessToken(ctx *gin.Context) {
	message := "could not refresh access token"

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

	context := ctx.Request.Context()

	user, err := userHandler.userService.UserGetById(context, fmt.Sprint(userID))

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": "the user is belonged to this token no longer exists "})
	}

	accessToken, err := utils.CreateToken(config.AccessTokenExpiresIn, user.UserID, config.AccessTokenPrivateKey)

	if err != nil {
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.SetCookie("access_token", accessToken, config.AccessTokenMaxAge*60, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "true", config.AccessTokenMaxAge*60, "/", "localhost", false, false)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "access_token": accessToken})
}

func (userHandler *UserHandler) LogoutUser(ctx *gin.Context) {
	ctx.SetCookie("access_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("refresh_token", "", -1, "/", "localhost", false, true)
	ctx.SetCookie("logged_in", "", -1, "/", "localhost", false, true)

	ctx.JSON(http.StatusOK, gin.H{"status": "success"})
}

func (userHandler *UserHandler) GetMe(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(*models.UserFullResponse)

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": gin.H{"user": models.FilteredResponse(currentUser)}})
}
