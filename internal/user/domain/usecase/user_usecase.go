package usecase

import (
	"context"
	"time"

	"github.com/thanhpk/randstr"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	logger "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logger"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location                   = "internal.user.domain.usecase."
	verificationCodeLength int = 20
	resetTokenLength       int = 20
	userRole                   = "user"
)

type UserUseCaseV1 struct {
	userRepository user.UserRepository
}

func NewUserUseCaseV1(userRepository user.UserRepository) user.UserUseCase {
	return UserUseCaseV1{userRepository: userRepository}
}

func (userUseCaseV1 UserUseCaseV1) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	fetchedUsers := userUseCaseV1.userRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		return commonModel.NewResultOnFailure[userModel.Users](domainError.HandleError(fetchedUsers.Error))
	}

	return fetchedUsers
}

func (userUseCaseV1 UserUseCaseV1) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	fetchedUser := userUseCaseV1.userRepository.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

func (userUseCaseV1 UserUseCaseV1) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	validateEmailError := checkEmail(location+"GetUserByEmail", email)
	if validator.IsError(validateEmailError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(validateEmailError))
	}

	fetchedUser := userUseCaseV1.userRepository.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

func (userUseCaseV1 UserUseCaseV1) Register(ctx context.Context, userCreateData userModel.UserCreate) commonModel.Result[userModel.User] {
	userCreate := validateUserCreate(userCreateData)
	if validator.IsError(userCreate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userCreate.Error))
	}

	checkEmailDuplicateError := userUseCaseV1.userRepository.CheckEmailDuplicate(ctx, userCreate.Data.Email)
	if validator.IsError(checkEmailDuplicateError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(checkEmailDuplicateError))
	}

	token := randstr.String(verificationCodeLength)
	token = commonUtility.Encode(token)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = token
	currentTime := time.Now()
	userCreate.Data.CreatedAt = currentTime
	userCreate.Data.UpdatedAt = currentTime

	createdUser := userUseCaseV1.userRepository.Register(ctx, userCreate.Data)
	if validator.IsError(createdUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(createdUser.Error))
	}

	emailData := prepareEmailDataForRegistration(createdUser.Data.Name, token)
	sendEmailError := userUseCaseV1.userRepository.SendEmail(createdUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(sendEmailError))
	}

	return createdUser
}

func (userUseCaseV1 UserUseCaseV1) UpdateCurrentUser(ctx context.Context, userUpdateData userModel.UserUpdate) commonModel.Result[userModel.User] {
	userUpdate := validateUserUpdate(userUpdateData)
	if validator.IsError(userUpdate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userUpdate.Error))
	}

	userUpdate.Data.UpdatedAt = time.Now()
	updatedUser := userUseCaseV1.userRepository.UpdateCurrentUser(ctx, userUpdate.Data)
	if validator.IsError(updatedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(updatedUser.Error))
	}

	return updatedUser
}

func (userUseCaseV1 UserUseCaseV1) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCaseV1.userRepository.DeleteUserById(ctx, userID)
	if validator.IsError(deletedUser) {
		return domainError.HandleError(deletedUser)
	}

	return nil
}

func (userUseCaseV1 UserUseCaseV1) Login(ctx context.Context, userLoginData userModel.UserLogin) commonModel.Result[userModel.UserToken] {
	userLogin := validateUserLogin(userLoginData)
	if validator.IsError(userLogin.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userLogin.Error))
	}

	fetchedUser := userUseCaseV1.userRepository.GetUserByEmail(ctx, userLogin.Data.Email)
	checkPasswordsError := checkPasswords(location+"Login", fetchedUser.Data.Password, userLoginData.Password)
	if validator.IsError(checkPasswordsError) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(checkPasswordsError))
	}

	userTokenPayload := domainModel.NewUserTokenPayload(fetchedUser.Data.ID, fetchedUser.Data.Role)
	userToken := generateToken(userTokenPayload)
	if validator.IsError(userToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userToken.Error))
	}

	return userToken
}

func (userUseCaseV1 UserUseCaseV1) RefreshAccessToken(ctx context.Context, user userModel.User) commonModel.Result[userModel.UserToken] {
	userTokenPayload := domainModel.NewUserTokenPayload(user.ID, user.Role)
	userToken := generateToken(userTokenPayload)
	if validator.IsError(userToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userToken.Error))
	}

	return userToken
}

func (userUseCaseV1 UserUseCaseV1) ForgottenPassword(ctx context.Context, userForgottenPasswordData userModel.UserForgottenPassword) error {
	userForgottenPassword := validateUserForgottenPassword(userForgottenPasswordData)
	if validator.IsError(userForgottenPassword.Error) {
		return domainError.HandleError(userForgottenPassword.Error)
	}

	fetchedUser := userUseCaseV1.GetUserByEmail(ctx, userForgottenPassword.Data.Email)
	if validator.IsError(fetchedUser.Error) {
		return domainError.HandleError(fetchedUser.Error)
	}

	token := randstr.String(resetTokenLength)
	encodedToken := commonUtility.Encode(token)
	tokenExpirationTime := time.Now().Add(constants.PasswordResetTokenExpirationTime)
	userForgottenPassword.Data.ResetToken = token
	userForgottenPassword.Data.ResetExpiry = tokenExpirationTime

	updatedUserPasswordError := userUseCaseV1.userRepository.ForgottenPassword(ctx, userForgottenPassword.Data)
	if validator.IsError(updatedUserPasswordError) {
		return domainError.HandleError(updatedUserPasswordError)
	}

	emailData := prepareEmailDataForForgottenPassword(fetchedUser.Data.Name, encodedToken)
	sendEmailError := userUseCaseV1.userRepository.SendEmail(fetchedUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return domainError.HandleError(sendEmailError)
	}

	return nil
}

func (userUseCaseV1 UserUseCaseV1) ResetUserPassword(ctx context.Context, userResetPasswordData userModel.UserResetPassword) error {
	token := commonUtility.Decode(location+"ResetUserPassword", userResetPasswordData.ResetToken)
	if validator.IsError(token.Error) {
		return domainError.HandleError(token.Error)
	}

	userResetPasswordData.ResetToken = token.Data
	userResetPassword := validateUserResetPassword(userResetPasswordData)
	if validator.IsError(userResetPassword.Error) {
		return domainError.HandleError(userResetPassword.Error)
	}

	fetchedResetExpiry := userUseCaseV1.userRepository.GetResetExpiry(ctx, token.Data)
	if validator.IsError(fetchedResetExpiry.Error) {
		return domainError.HandleError(fetchedResetExpiry.Error)
	}
	if validator.IsTimeNotValid(fetchedResetExpiry.Data.ResetExpiry) {
		timeExpiredError := domainError.NewTimeExpiredError(location+"ResetUserPassword", constants.TimeExpiredErrorNotification)
		logger.Logger(timeExpiredError)
		return domainError.HandleError(timeExpiredError)
	}

	resetUserPasswordError := userUseCaseV1.userRepository.ResetUserPassword(ctx, userResetPassword.Data)
	if validator.IsError(resetUserPasswordError) {
		return domainError.HandleError(resetUserPasswordError)
	}

	return nil
}

func prepareEmailData(userName, tokenValue, subject, url, templateName, templatePath string) userModel.EmailData {
	emailConfig := config.GetEmailConfig()
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(
		emailConfig.ClientOriginUrl+url+tokenValue,
		templateName,
		templatePath,
		userFirstName,
		subject,
	)

	return emailData
}

func prepareEmailDataForRegistration(userName, tokenValue string) userModel.EmailData {
	emailConfig := config.GetEmailConfig()
	return prepareEmailData(
		userName,
		tokenValue,
		constants.EmailConfirmationSubject,
		constants.EmailConfirmationUrl,
		emailConfig.UserConfirmationTemplateName,
		emailConfig.UserConfirmationTemplatePath,
	)
}

func prepareEmailDataForForgottenPassword(userName, tokenValue string) userModel.EmailData {
	emailConfig := config.GetEmailConfig()
	return prepareEmailData(
		userName,
		tokenValue,
		constants.ForgottenPasswordSubject,
		constants.ForgottenPasswordUrl,
		emailConfig.ForgottenPasswordTemplateName,
		emailConfig.ForgottenPasswordTemplatePath,
	)
}

func generateToken(userTokenPayload domainModel.UserTokenPayload) commonModel.Result[userModel.UserToken] {
	var userToken userModel.UserToken
	accessTokenConfig := config.GetAccessConfig()
	refreshTokenConfig := config.GetRefreshConfig()

	accessToken := domainUtility.GenerateJWTToken(location+".generateToken.accessToken", accessTokenConfig.PrivateKey, accessTokenConfig.ExpiredIn, userTokenPayload)
	if validator.IsError(accessToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](accessToken.Error)
	}

	refreshToken := domainUtility.GenerateJWTToken(location+".generateToken.refreshToken", refreshTokenConfig.PrivateKey, refreshTokenConfig.ExpiredIn, userTokenPayload)
	if validator.IsError(refreshToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](refreshToken.Error)
	}

	userToken.AccessToken = accessToken.Data
	userToken.RefreshToken = refreshToken.Data
	return commonModel.NewResultOnSuccess[userModel.UserToken](userToken)
}
