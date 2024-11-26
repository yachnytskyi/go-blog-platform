package usecase

import (
	"context"
	"time"

	"github.com/thanhpk/randstr"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/interfaces"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domain "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location                   = "internal.user.domain.usecase."
	verificationCodeLength int = 20
	resetTokenLength       int = 20
	userRole                   = "user"
)

type UserUseCase struct {
	Config         *config.ApplicationConfig
	Logger         interfaces.Logger
	Email          interfaces.Email
	UserRepository interfaces.UserRepository
}

func NewUserUseCase(config *config.ApplicationConfig, logger interfaces.Logger, email interfaces.Email, userRepository interfaces.UserRepository) UserUseCase {
	return UserUseCase{
		Config:         config,
		Logger:         logger,
		Email:          email,
		UserRepository: userRepository,
	}
}

func (userUseCase UserUseCase) GetAllUsers(ctx context.Context, paginationQuery common.PaginationQuery) common.Result[user.Users] {
	fetchedUsers := userUseCase.UserRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		return common.NewResultOnFailure[user.Users](domain.HandleError(fetchedUsers.Error))
	}

	return fetchedUsers
}

func (userUseCase UserUseCase) GetUserById(ctx context.Context, userID string) common.Result[user.User] {
	fetchedUser := userUseCase.UserRepository.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		return common.NewResultOnFailure[user.User](domain.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

func (userUseCase UserUseCase) GetUserByEmail(ctx context.Context, email string) common.Result[user.User] {
	validateEmailError := checkEmail(userUseCase.Logger, location+"GetUserByEmail", email)
	if validator.IsError(validateEmailError) {
		return common.NewResultOnFailure[user.User](domain.HandleError(validateEmailError))
	}

	fetchedUser := userUseCase.UserRepository.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		return common.NewResultOnFailure[user.User](domain.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

func (userUseCase UserUseCase) Register(ctx context.Context, userCreateData user.UserCreate) common.Result[user.User] {
	userCreate := validateUserCreate(userUseCase.Logger, userCreateData)
	if validator.IsError(userCreate.Error) {
		return common.NewResultOnFailure[user.User](domain.HandleError(userCreate.Error))
	}

	checkEmailDuplicateError := userUseCase.UserRepository.CheckEmailDuplicate(ctx, userCreate.Data.Email)
	if validator.IsError(checkEmailDuplicateError) {
		return common.NewResultOnFailure[user.User](domain.HandleError(checkEmailDuplicateError))
	}

	token := randstr.String(verificationCodeLength)
	encodedToken := utility.Encode(token)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = token

	createdUser := userUseCase.UserRepository.Register(ctx, userCreate.Data)
	if validator.IsError(createdUser.Error) {
		return common.NewResultOnFailure[user.User](domain.HandleError(createdUser.Error))
	}

	emailData := prepareEmailDataForRegistration(userUseCase.Config, createdUser.Data, encodedToken)
	sendEmailError := userUseCase.Email.SendEmail(userUseCase.Config, userUseCase.Logger, location+"Register", createdUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return common.NewResultOnFailure[user.User](domain.HandleError(sendEmailError))
	}

	return createdUser
}

func (userUseCase UserUseCase) UpdateCurrentUser(ctx context.Context, userUpdateData user.UserUpdate) common.Result[user.User] {
	userUpdate := validateUserUpdate(userUseCase.Logger, userUpdateData)
	if validator.IsError(userUpdate.Error) {
		return common.NewResultOnFailure[user.User](domain.HandleError(userUpdate.Error))
	}

	updatedUser := userUseCase.UserRepository.UpdateCurrentUser(ctx, userUpdate.Data)
	if validator.IsError(updatedUser.Error) {
		return common.NewResultOnFailure[user.User](domain.HandleError(updatedUser.Error))
	}

	return updatedUser
}

func (userUseCase UserUseCase) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCase.UserRepository.DeleteUserById(ctx, userID)
	if validator.IsError(deletedUser) {
		return domain.HandleError(deletedUser)
	}

	return nil
}

func (userUseCase UserUseCase) Login(ctx context.Context, userLoginData user.UserLogin) common.Result[user.UserToken] {
	userLogin := validateUserLogin(userUseCase.Logger, userLoginData)
	if validator.IsError(userLogin.Error) {
		return common.NewResultOnFailure[user.UserToken](domain.HandleError(userLogin.Error))
	}

	fetchedUser := userUseCase.UserRepository.GetUserByEmail(ctx, userLogin.Data.Email)
	checkPasswordsError := checkPasswords(userUseCase.Logger, location+"Login", fetchedUser.Data.Password, userLoginData.Password)
	if validator.IsError(checkPasswordsError) {
		return common.NewResultOnFailure[user.UserToken](domain.HandleError(checkPasswordsError))
	}

	userTokenPayload := model.NewUserTokenPayload(fetchedUser.Data.ID, fetchedUser.Data.Role)
	userToken := generateToken(userUseCase.Config, userUseCase.Logger, location+"Login", userTokenPayload)
	if validator.IsError(userToken.Error) {
		return common.NewResultOnFailure[user.UserToken](domain.HandleError(userToken.Error))
	}

	return userToken
}

func (userUseCase UserUseCase) RefreshAccessToken(ctx context.Context, userData user.User) common.Result[user.UserToken] {
	userTokenPayload := model.NewUserTokenPayload(userData.ID, userData.Role)
	userToken := generateToken(userUseCase.Config, userUseCase.Logger, location+"RefreshAccessToken", userTokenPayload)
	if validator.IsError(userToken.Error) {
		return common.NewResultOnFailure[user.UserToken](domain.HandleError(userToken.Error))
	}

	return userToken
}

func (userUseCase UserUseCase) ForgottenPassword(ctx context.Context, userForgottenPasswordData user.UserForgottenPassword) error {
	userForgottenPassword := validateUserForgottenPassword(userUseCase.Logger, userForgottenPasswordData)
	if validator.IsError(userForgottenPassword.Error) {
		return domain.HandleError(userForgottenPassword.Error)
	}

	fetchedUser := userUseCase.GetUserByEmail(ctx, userForgottenPassword.Data.Email)
	if validator.IsError(fetchedUser.Error) {
		return domain.HandleError(fetchedUser.Error)
	}

	token := randstr.String(resetTokenLength)
	encodedToken := utility.Encode(token)
	tokenExpirationTime := time.Now().Add(constants.PasswordResetTokenExpirationTime)
	userForgottenPassword.Data.ResetToken = token
	userForgottenPassword.Data.ResetExpiry = tokenExpirationTime

	updatedUserPasswordError := userUseCase.UserRepository.ForgottenPassword(ctx, userForgottenPassword.Data)
	if validator.IsError(updatedUserPasswordError) {
		return domain.HandleError(updatedUserPasswordError)
	}

	emailData := prepareEmailDataForForgottenPassword(userUseCase.Config, fetchedUser.Data, encodedToken)
	sendEmailError := userUseCase.Email.SendEmail(userUseCase.Config, userUseCase.Logger, location+"ForgottenPassword", fetchedUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return domain.HandleError(sendEmailError)
	}

	return nil
}

func (userUseCase UserUseCase) ResetUserPassword(ctx context.Context, userResetPasswordData user.UserResetPassword) error {
	token := utility.Decode(userUseCase.Logger, location+"ResetUserPassword", userResetPasswordData.ResetToken)
	if validator.IsError(token.Error) {
		return domain.HandleError(token.Error)
	}

	userResetPasswordData.ResetToken = token.Data
	userResetPassword := validateUserResetPassword(userUseCase.Logger, userResetPasswordData)
	if validator.IsError(userResetPassword.Error) {
		return domain.HandleError(userResetPassword.Error)
	}

	fetchedResetExpiry := userUseCase.UserRepository.GetResetExpiry(ctx, token.Data)
	if validator.IsError(fetchedResetExpiry.Error) {
		return domain.HandleError(fetchedResetExpiry.Error)
	}
	if validator.IsTimeNotValid(fetchedResetExpiry.Data.ResetExpiry) {
		timeExpiredError := domain.NewTimeExpiredError(location+"ResetUserPassword.IsTimeNotValid", constants.TimeExpiredErrorNotification)
		userUseCase.Logger.Error(timeExpiredError)
		return domain.HandleError(timeExpiredError)
	}

	resetUserPasswordError := userUseCase.UserRepository.ResetUserPassword(ctx, userResetPassword.Data)
	if validator.IsError(resetUserPasswordError) {
		return domain.HandleError(resetUserPasswordError)
	}

	return nil
}

func prepareEmailData(config *config.ApplicationConfig, user user.User, tokenValue, subject, url, templateName, templatePath string) interfaces.EmailData {
	userFirstName := domainUtility.UserFirstName(user.Username)
	emailData := interfaces.NewEmailData(
		user.Email,
		config.Email.ClientOriginUrl+url+tokenValue,
		templateName,
		templatePath,
		userFirstName,
		subject,
	)

	return emailData
}

func prepareEmailDataForRegistration(config *config.ApplicationConfig, user user.User, tokenValue string) interfaces.EmailData {
	return prepareEmailData(
		config,
		user,
		tokenValue,
		constants.EmailConfirmationSubject,
		constants.EmailConfirmationUrl,
		config.Email.UserConfirmationTemplateName,
		config.Email.UserConfirmationTemplatePath,
	)
}

func prepareEmailDataForForgottenPassword(config *config.ApplicationConfig, user user.User, tokenValue string) interfaces.EmailData {
	return prepareEmailData(
		config,
		user,
		tokenValue,
		constants.ForgottenPasswordSubject,
		constants.ForgottenPasswordUrl,
		config.Email.ForgottenPasswordTemplateName,
		config.Email.ForgottenPasswordTemplatePath,
	)
}

func generateToken(config *config.ApplicationConfig, logger interfaces.Logger, location string, userTokenPayload model.UserTokenPayload) common.Result[user.UserToken] {
	accessToken := domainUtility.GenerateJWTToken(
		logger,
		location+".generateToken.accessToken",
		config.AccessToken.PrivateKey,
		config.AccessToken.ExpiredIn,
		userTokenPayload,
	)
	if validator.IsError(accessToken.Error) {
		return common.NewResultOnFailure[user.UserToken](accessToken.Error)
	}

	refreshToken := domainUtility.GenerateJWTToken(
		logger,
		location+".generateToken.refreshToken",
		config.RefreshToken.PrivateKey,
		config.RefreshToken.ExpiredIn,
		userTokenPayload,
	)
	if validator.IsError(refreshToken.Error) {
		return common.NewResultOnFailure[user.UserToken](refreshToken.Error)
	}

	var userToken user.UserToken
	userToken.AccessToken = accessToken.Data
	userToken.RefreshToken = refreshToken.Data
	return common.NewResultOnSuccess[user.UserToken](userToken)
}
