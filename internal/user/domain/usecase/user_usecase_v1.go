package usecase

import (
	"context"
	"time"

	"github.com/thanhpk/randstr"
	constants "github.com/yachnytskyi/golang-mongo-grpc/config/constants"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	config "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/factory/config/model"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/dependency/model"
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
	Config         model.Config
	Logger         model.Logger
	UserRepository user.UserRepository
}

func NewUserUseCaseV1(config model.Config, logger model.Logger, userRepository user.UserRepository) user.UserUseCase {
	return &UserUseCaseV1{
		Config:         config,
		Logger:         logger,
		UserRepository: userRepository,
	}
}

func (userUseCaseV1 *UserUseCaseV1) GetAllUsers(ctx context.Context, paginationQuery common.PaginationQuery) common.Result[userModel.Users] {
	fetchedUsers := userUseCaseV1.UserRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		return common.NewResultOnFailure[userModel.Users](domainError.HandleError(fetchedUsers.Error))
	}

	return fetchedUsers
}

func (userUseCaseV1 *UserUseCaseV1) GetUserById(ctx context.Context, userID string) common.Result[userModel.User] {
	fetchedUser := userUseCaseV1.UserRepository.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

func (userUseCaseV1 *UserUseCaseV1) GetUserByEmail(ctx context.Context, email string) common.Result[userModel.User] {
	validateEmailError := checkEmail(userUseCaseV1.Logger, location+"GetUserByEmail", email)
	if validator.IsError(validateEmailError) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(validateEmailError))
	}

	fetchedUser := userUseCaseV1.UserRepository.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

func (userUseCaseV1 *UserUseCaseV1) Register(ctx context.Context, userCreateData userModel.UserCreate) common.Result[userModel.User] {
	userCreate := validateUserCreate(userUseCaseV1.Logger, userCreateData)
	if validator.IsError(userCreate.Error) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(userCreate.Error))
	}

	checkEmailDuplicateError := userUseCaseV1.UserRepository.CheckEmailDuplicate(ctx, userCreate.Data.Email)
	if validator.IsError(checkEmailDuplicateError) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(checkEmailDuplicateError))
	}

	token := randstr.String(verificationCodeLength)
	token = utility.Encode(token)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = token
	currentTime := time.Now()
	userCreate.Data.CreatedAt = currentTime
	userCreate.Data.UpdatedAt = currentTime

	createdUser := userUseCaseV1.UserRepository.Register(ctx, userCreate.Data)
	if validator.IsError(createdUser.Error) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(createdUser.Error))
	}

	emailData := prepareEmailDataForRegistration(userUseCaseV1.Config, createdUser.Data.Name, token)
	sendEmailError := userUseCaseV1.UserRepository.SendEmail(createdUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(sendEmailError))
	}

	return createdUser
}

func (userUseCaseV1 *UserUseCaseV1) UpdateCurrentUser(ctx context.Context, userUpdateData userModel.UserUpdate) common.Result[userModel.User] {
	userUpdate := validateUserUpdate(userUseCaseV1.Logger, userUpdateData)
	if validator.IsError(userUpdate.Error) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(userUpdate.Error))
	}

	userUpdate.Data.UpdatedAt = time.Now()
	updatedUser := userUseCaseV1.UserRepository.UpdateCurrentUser(ctx, userUpdate.Data)
	if validator.IsError(updatedUser.Error) {
		return common.NewResultOnFailure[userModel.User](domainError.HandleError(updatedUser.Error))
	}

	return updatedUser
}

func (userUseCaseV1 *UserUseCaseV1) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCaseV1.UserRepository.DeleteUserById(ctx, userID)
	if validator.IsError(deletedUser) {
		return domainError.HandleError(deletedUser)
	}

	return nil
}

func (userUseCaseV1 *UserUseCaseV1) Login(ctx context.Context, userLoginData userModel.UserLogin) common.Result[userModel.UserToken] {
	userLogin := validateUserLogin(userUseCaseV1.Logger, userLoginData)
	if validator.IsError(userLogin.Error) {
		return common.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userLogin.Error))
	}

	fetchedUser := userUseCaseV1.UserRepository.GetUserByEmail(ctx, userLogin.Data.Email)
	checkPasswordsError := checkPasswords(userUseCaseV1.Logger, location+"Login", fetchedUser.Data.Password, userLoginData.Password)
	if validator.IsError(checkPasswordsError) {
		return common.NewResultOnFailure[userModel.UserToken](domainError.HandleError(checkPasswordsError))
	}

	userTokenPayload := domainModel.NewUserTokenPayload(fetchedUser.Data.ID, fetchedUser.Data.Role)
	userToken := generateToken(userUseCaseV1.Config, userUseCaseV1.Logger, userTokenPayload)
	if validator.IsError(userToken.Error) {
		return common.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userToken.Error))
	}

	return userToken
}

func (userUseCaseV1 *UserUseCaseV1) RefreshAccessToken(ctx context.Context, user userModel.User) common.Result[userModel.UserToken] {
	userTokenPayload := domainModel.NewUserTokenPayload(user.ID, user.Role)
	userToken := generateToken(userUseCaseV1.Config, userUseCaseV1.Logger, userTokenPayload)
	if validator.IsError(userToken.Error) {
		return common.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userToken.Error))
	}

	return userToken
}

func (userUseCaseV1 *UserUseCaseV1) ForgottenPassword(ctx context.Context, userForgottenPasswordData userModel.UserForgottenPassword) error {
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

	emailData := prepareEmailDataForForgottenPassword(userUseCaseV1.Config, fetchedUser.Data.Name, encodedToken)
	sendEmailError := userUseCaseV1.UserRepository.SendEmail(fetchedUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return domainError.HandleError(sendEmailError)
	}

	return nil
}

func (userUseCaseV1 *UserUseCaseV1) ResetUserPassword(ctx context.Context, userResetPasswordData userModel.UserResetPassword) error {
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
		timeExpiredError := domainError.NewTimeExpiredError(location+"ResetUserPassword", constants.TimeExpiredErrorNotification)
		userUseCaseV1.Logger.Error(timeExpiredError)
		return domainError.HandleError(timeExpiredError)
	}

	resetUserPasswordError := userUseCaseV1.UserRepository.ResetUserPassword(ctx, userResetPassword.Data)
	if validator.IsError(resetUserPasswordError) {
		return domainError.HandleError(resetUserPasswordError)
	}

	return nil
}

func prepareEmailData(config *config.ApplicationConfig, userName, tokenValue, subject, url, templateName, templatePath string) userModel.EmailData {
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(
		config.Email.ClientOriginUrl+url+tokenValue,
		templateName,
		templatePath,
		userFirstName,
		subject,
	)

	return emailData
}

func prepareEmailDataForRegistration(configInstance model.Config, userName, tokenValue string) userModel.EmailData {
	config := configInstance.GetConfig()
	return prepareEmailData(
		config,
		userName,
		tokenValue,
		constants.EmailConfirmationSubject,
		constants.EmailConfirmationUrl,
		config.Email.UserConfirmationTemplateName,
		config.Email.UserConfirmationTemplatePath,
	)
}

func prepareEmailDataForForgottenPassword(configInstance model.Config, userName, tokenValue string) userModel.EmailData {
	config := configInstance.GetConfig()
	return prepareEmailData(
		config,
		userName,
		tokenValue,
		constants.ForgottenPasswordSubject,
		constants.ForgottenPasswordUrl,
		config.Email.ForgottenPasswordTemplateName,
		config.Email.ForgottenPasswordTemplatePath,
	)
}

func generateToken(config model.Config, logger model.Logger, userTokenPayload domainModel.UserTokenPayload) common.Result[userModel.UserToken] {
	configInstance := config.GetConfig()

	accessToken := domainUtility.GenerateJWTToken(
		logger,
		location+".generateToken.accessToken",
		configInstance.AccessToken.PrivateKey,
		configInstance.AccessToken.ExpiredIn,
		userTokenPayload,
	)
	if validator.IsError(accessToken.Error) {
		return common.NewResultOnFailure[userModel.UserToken](accessToken.Error)
	}

	refreshToken := domainUtility.GenerateJWTToken(
		logger,
		location+".generateToken.refreshToken",
		configInstance.RefreshToken.PrivateKey,
		configInstance.RefreshToken.ExpiredIn,
		userTokenPayload,
	)
	if validator.IsError(refreshToken.Error) {
		return common.NewResultOnFailure[userModel.UserToken](refreshToken.Error)
	}

	var userToken userModel.UserToken
	userToken.AccessToken = accessToken.Data
	userToken.RefreshToken = refreshToken.Data
	return common.NewResultOnSuccess[userModel.UserToken](userToken)
}
