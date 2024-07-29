package usecase

import (
	"context"
	"time"

	"github.com/thanhpk/randstr"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	interfaces "github.com/yachnytskyi/golang-mongo-grpc/internal/common/interfaces"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location                   = "internal.user.domain.usecase."
	verificationCodeLength int = 20
	resetTokenLength       int = 20
	userRole                   = "user"
)

type UserUseCaseV1 struct {
	Config         interfaces.Config
	Logger         interfaces.Logger
	Email          interfaces.Email
	UserRepository interfaces.UserRepository
}

func NewUserUseCaseV1(config interfaces.Config, logger interfaces.Logger, email interfaces.Email, userRepository interfaces.UserRepository) UserUseCaseV1 {
	return UserUseCaseV1{
		Config:         config,
		Logger:         logger,
		Email:          email,
		UserRepository: userRepository,
	}
}

func (userUseCaseV1 UserUseCaseV1) GetAllUsers(ctx context.Context, paginationQuery common.PaginationQuery) common.Result[user.Users] {
	fetchedUsers := userUseCaseV1.UserRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		return common.NewResultOnFailure[user.Users](domainError.HandleError(fetchedUsers.Error))
	}

	return fetchedUsers
}

func (userUseCaseV1 UserUseCaseV1) GetUserById(ctx context.Context, userID string) common.Result[user.User] {
	fetchedUser := userUseCaseV1.UserRepository.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

func (userUseCaseV1 UserUseCaseV1) GetUserByEmail(ctx context.Context, email string) common.Result[user.User] {
	validateEmailError := checkEmail(userUseCaseV1.Logger, location+"GetUserByEmail", email)
	if validator.IsError(validateEmailError) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(validateEmailError))
	}

	fetchedUser := userUseCaseV1.UserRepository.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

func (userUseCaseV1 UserUseCaseV1) Register(ctx context.Context, userCreateData user.UserCreate) common.Result[user.User] {
	userCreate := validateUserCreate(userUseCaseV1.Logger, userCreateData)
	if validator.IsError(userCreate.Error) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(userCreate.Error))
	}

	checkEmailDuplicateError := userUseCaseV1.UserRepository.CheckEmailDuplicate(ctx, userCreate.Data.Email)
	if validator.IsError(checkEmailDuplicateError) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(checkEmailDuplicateError))
	}

	token := randstr.String(verificationCodeLength)
	encodedToken := utility.Encode(token)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = token
	currentTime := time.Now()
	userCreate.Data.CreatedAt = currentTime
	userCreate.Data.UpdatedAt = currentTime

	createdUser := userUseCaseV1.UserRepository.Register(ctx, userCreate.Data)
	if validator.IsError(createdUser.Error) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(createdUser.Error))
	}

	emailData := prepareEmailDataForRegistration(userUseCaseV1.Config, createdUser.Data, encodedToken)
	sendEmailError := userUseCaseV1.Email.SendEmail(userUseCaseV1.Config, userUseCaseV1.Logger, location+"Register", createdUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(sendEmailError))
	}

	return createdUser
}

func (userUseCaseV1 UserUseCaseV1) UpdateCurrentUser(ctx context.Context, userUpdateData user.UserUpdate) common.Result[user.User] {
	userUpdate := validateUserUpdate(userUseCaseV1.Logger, userUpdateData)
	if validator.IsError(userUpdate.Error) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(userUpdate.Error))
	}

	userUpdate.Data.UpdatedAt = time.Now()
	updatedUser := userUseCaseV1.UserRepository.UpdateCurrentUser(ctx, userUpdate.Data)
	if validator.IsError(updatedUser.Error) {
		return common.NewResultOnFailure[user.User](domainError.HandleError(updatedUser.Error))
	}

	return updatedUser
}

func (userUseCaseV1 UserUseCaseV1) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCaseV1.UserRepository.DeleteUserById(ctx, userID)
	if validator.IsError(deletedUser) {
		return domainError.HandleError(deletedUser)
	}

	return nil
}

func (userUseCaseV1 UserUseCaseV1) Login(ctx context.Context, userLoginData user.UserLogin) common.Result[user.UserToken] {
	userLogin := validateUserLogin(userUseCaseV1.Logger, userLoginData)
	if validator.IsError(userLogin.Error) {
		return common.NewResultOnFailure[user.UserToken](domainError.HandleError(userLogin.Error))
	}

	fetchedUser := userUseCaseV1.UserRepository.GetUserByEmail(ctx, userLogin.Data.Email)
	checkPasswordsError := checkPasswords(userUseCaseV1.Logger, location+"Login", fetchedUser.Data.Password, userLoginData.Password)
	if validator.IsError(checkPasswordsError) {
		return common.NewResultOnFailure[user.UserToken](domainError.HandleError(checkPasswordsError))
	}

	userTokenPayload := domainModel.NewUserTokenPayload(fetchedUser.Data.ID, fetchedUser.Data.Role)
	userToken := generateToken(userUseCaseV1.Config, userUseCaseV1.Logger, userTokenPayload)
	if validator.IsError(userToken.Error) {
		return common.NewResultOnFailure[user.UserToken](domainError.HandleError(userToken.Error))
	}

	return userToken
}

func (userUseCaseV1 UserUseCaseV1) RefreshAccessToken(ctx context.Context, userData user.User) common.Result[user.UserToken] {
	userTokenPayload := domainModel.NewUserTokenPayload(userData.ID, userData.Role)
	userToken := generateToken(userUseCaseV1.Config, userUseCaseV1.Logger, userTokenPayload)
	if validator.IsError(userToken.Error) {
		return common.NewResultOnFailure[user.UserToken](domainError.HandleError(userToken.Error))
	}

	return userToken
}

func (userUseCaseV1 UserUseCaseV1) ForgottenPassword(ctx context.Context, userForgottenPasswordData user.UserForgottenPassword) error {
	userForgottenPassword := validateUserForgottenPassword(userUseCaseV1.Logger, userForgottenPasswordData)
	if validator.IsError(userForgottenPassword.Error) {
		return domainError.HandleError(userForgottenPassword.Error)
	}

	fetchedUser := userUseCaseV1.GetUserByEmail(ctx, userForgottenPassword.Data.Email)
	if validator.IsError(fetchedUser.Error) {
		return domainError.HandleError(fetchedUser.Error)
	}

	token := randstr.String(resetTokenLength)
	encodedToken := utility.Encode(token)
	tokenExpirationTime := time.Now().Add(constants.PasswordResetTokenExpirationTime)
	userForgottenPassword.Data.ResetToken = token
	userForgottenPassword.Data.ResetExpiry = tokenExpirationTime

	updatedUserPasswordError := userUseCaseV1.UserRepository.ForgottenPassword(ctx, userForgottenPassword.Data)
	if validator.IsError(updatedUserPasswordError) {
		return domainError.HandleError(updatedUserPasswordError)
	}

	emailData := prepareEmailDataForForgottenPassword(userUseCaseV1.Config, fetchedUser.Data, encodedToken)
	sendEmailError := userUseCaseV1.Email.SendEmail(userUseCaseV1.Config, userUseCaseV1.Logger, location+"ForgottenPassword", fetchedUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return domainError.HandleError(sendEmailError)
	}

	return nil
}

func (userUseCaseV1 UserUseCaseV1) ResetUserPassword(ctx context.Context, userResetPasswordData user.UserResetPassword) error {
	token := utility.Decode(userUseCaseV1.Logger, location+"ResetUserPassword", userResetPasswordData.ResetToken)
	if validator.IsError(token.Error) {
		return domainError.HandleError(token.Error)
	}

	userResetPasswordData.ResetToken = token.Data
	userResetPassword := validateUserResetPassword(userUseCaseV1.Logger, userResetPasswordData)
	if validator.IsError(userResetPassword.Error) {
		return domainError.HandleError(userResetPassword.Error)
	}

	fetchedResetExpiry := userUseCaseV1.UserRepository.GetResetExpiry(ctx, token.Data)
	if validator.IsError(fetchedResetExpiry.Error) {
		return domainError.HandleError(fetchedResetExpiry.Error)
	}
	if validator.IsTimeNotValid(fetchedResetExpiry.Data.ResetExpiry) {
		timeExpiredError := domainError.NewTimeExpiredError(location+"ResetUserPassword.IsTimeNotValid", constants.TimeExpiredErrorNotification)
		userUseCaseV1.Logger.Error(timeExpiredError)
		return domainError.HandleError(timeExpiredError)
	}

	resetUserPasswordError := userUseCaseV1.UserRepository.ResetUserPassword(ctx, userResetPassword.Data)
	if validator.IsError(resetUserPasswordError) {
		return domainError.HandleError(resetUserPasswordError)
	}

	return nil
}

func prepareEmailData(config *config.ApplicationConfig, user user.User, tokenValue, subject, url, templateName, templatePath string) interfaces.EmailData {
	userFirstName := domainUtility.UserFirstName(user.Name)
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

func prepareEmailDataForRegistration(configInstance interfaces.Config, user user.User, tokenValue string) interfaces.EmailData {
	config := configInstance.GetConfig()
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

func prepareEmailDataForForgottenPassword(configInstance interfaces.Config, user user.User, tokenValue string) interfaces.EmailData {
	config := configInstance.GetConfig()
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

func generateToken(config interfaces.Config, logger interfaces.Logger, userTokenPayload domainModel.UserTokenPayload) common.Result[user.UserToken] {
	configInstance := config.GetConfig()

	accessToken := domainUtility.GenerateJWTToken(
		logger,
		location+".generateToken.accessToken",
		configInstance.AccessToken.PrivateKey,
		configInstance.AccessToken.ExpiredIn,
		userTokenPayload,
	)
	if validator.IsError(accessToken.Error) {
		return common.NewResultOnFailure[user.UserToken](accessToken.Error)
	}

	refreshToken := domainUtility.GenerateJWTToken(
		logger,
		location+".generateToken.refreshToken",
		configInstance.RefreshToken.PrivateKey,
		configInstance.RefreshToken.ExpiredIn,
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
