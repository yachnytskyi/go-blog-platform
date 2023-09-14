package usecase

import (
	"context"
	"time"

	"github.com/thanhpk/randstr"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	verificationCodeLength   int    = 20
	resetTokenLength         int    = 20
	emailConfirmationUrl     string = "users/verifyemail/"
	forgottenPasswordUrl     string = "users/reset-password/"
	emailConfirmationSubject string = "Your account verification code"
	forgottenPasswordSubject string = "Your password reset token (it is valid for 15 minutes)"
)

type UserUseCase struct {
	userRepository user.Repository
}

func NewUserUseCase(userRepository user.Repository) user.UseCase {
	return &UserUseCase{userRepository: userRepository}
}

func (userUseCase UserUseCase) GetAllUsers(ctx context.Context, page int, limit int) (userModel.Users, error) {
	fetchedUsers, getAllUsers := userUseCase.userRepository.GetAllUsers(ctx, page, limit)
	return fetchedUsers, getAllUsers
}

func (userUseCase UserUseCase) GetUserById(ctx context.Context, userID string) (userModel.User, error) {
	fetchedUser, getUserByIdError := userUseCase.userRepository.GetUserById(ctx, userID)
	return fetchedUser, getUserByIdError
}

func (userUseCase UserUseCase) GetUserByEmail(ctx context.Context, email string) (userModel.User, error) {
	fetchedUser, getUserByEmailError := userUseCase.userRepository.GetUserByEmail(ctx, email)
	return fetchedUser, getUserByEmailError
}

func (userUseCase UserUseCase) Register(ctx context.Context, userCreateData userModel.UserCreate) common.Result[userModel.User] {
	userCreate := validateUserCreate(userCreateData)
	if validator.IsErrorNotNil(userCreate.Error) {
		return common.NewResultWithError[userModel.User](domainError.HandleError(userCreate.Error))
	}
	checkEmailDublicateError := userUseCase.userRepository.CheckEmailDublicate(ctx, userCreate.Data.Email)
	if validator.IsErrorNotNil(checkEmailDublicateError) {
		checkEmailDublicateError = domainError.HandleError(checkEmailDublicateError)
		return common.NewResultWithError[userModel.User](checkEmailDublicateError)
	}
	userCreate.Data.Verified = true
	userCreate.Data.Role = "user"

	createdUser := userUseCase.userRepository.Register(ctx, userCreate.Data)
	if validator.IsErrorNotNil(createdUser.Error) {
		createdUser.Error = domainError.HandleError(createdUser.Error)
		return createdUser
	}

	// // Generate verification code.
	// tokenValue := randstr.String(verificationCodeLength)
	// encodedTokenValue := commonUtility.Encode(tokenValue)
	// _, userUpdateError := userUseCase.userRepository.UpdateNewRegisteredUserById(ctx, createdUser.Data.UserID, "verificationCode", encodedTokenValue)
	// if validator.IsErrorNotNil(userUpdateError) {
	// 	userUpdateError = domainError.HandleError(userUpdateError)
	// 	return userUpdateError
	// }
	// templateName := config.UserConfirmationEmailTemplateName
	// templatePath := config.UserConfirmationEmailTemplatePath

	// emailData, prepareEmailDataError := PrepareEmailData(ctx, createdUser.Data.Name, emailConfirmationUrl, emailConfirmationSubject, tokenValue, templateName, templatePath)
	// if validator.IsErrorNotNil(prepareEmailDataError) {
	// 	logging.Logger(prepareEmailDataError)
	// 	prepareEmailDataError = domainError.HandleError(prepareEmailDataError)
	// 	return prepareEmailDataError
	// }
	// sendEmailVerificationMessageError := userUseCase.userRepository.SendEmailVerificationMessage(ctx, createdUser.Data, emailData)
	// if validator.IsErrorNotNil(sendEmailVerificationMessageError) {
	// 	sendEmailVerificationMessageError = domainError.HandleError(sendEmailVerificationMessageError)
	// 	return sendEmailVerificationMessageError
	// }
	return createdUser
}

func (userUseCase UserUseCase) UpdateUserById(ctx context.Context, userID string, userUpdateData userModel.UserUpdate) (userModel.User, error) {
	userUpdate := validateUserUpdate(userUpdateData)
	if validator.IsErrorNotNil(userUpdate.Error) {
		validationErrors := domainError.HandleError(userUpdate.Error)
		return userModel.User{}, validationErrors
	}

	updatedUser, userUpdateError := userUseCase.userRepository.UpdateUserById(ctx, userID, userUpdate.Data)
	if validator.IsErrorNotNil(userUpdateError) {
		userUpdateError = domainError.HandleError(userUpdateError)
		return userModel.User{}, userUpdateError
	}
	return updatedUser, nil
}

func (userUseCase UserUseCase) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCase.userRepository.DeleteUserById(ctx, userID)
	return deletedUser
}

func (userUseCase UserUseCase) Login(ctx context.Context, userLoginData userModel.UserLogin) (string, error) {
	userLogin := validateUserLogin(userLoginData)
	if validator.IsErrorNotNil(userLogin.Error) {
		validationErrors := domainError.HandleError(userLogin.Error)
		return "", validationErrors
	}

	fetchedUser, getUserByEmailError := userUseCase.userRepository.GetUserByEmail(ctx, userLogin.Data.Email)
	if validator.IsErrorNotNil(getUserByEmailError) {
		return "", domainError.HandleError(domainError.NewErrorMessage(invalidEmailOrPassword))
	}
	if ArePasswordsNotEqual(fetchedUser.Password, userLoginData.Password) {
		return "", domainError.HandleError(domainError.NewErrorMessage(invalidEmailOrPassword))
	}
	return fetchedUser.UserID, nil
}

func (userUseCase UserUseCase) UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (userModel.User, error) {
	updatedUser, updateNewRegisteredUserById := userUseCase.userRepository.UpdateNewRegisteredUserById(ctx, userID, key, value)
	return updatedUser, updateNewRegisteredUserById
}

func (userUseCase UserUseCase) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {
	updatedUserError := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, email, firstKey, firstValue, secondKey, secondValue)
	if validator.IsErrorNotNil(updatedUserError) {
		updatedUserError = domainError.HandleError(updatedUserError)
		return updatedUserError
	}

	// Generate verification code.
	tokenValue := randstr.String(resetTokenLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	tokenExpirationTime := time.Now().Add(time.Minute * 15)

	// Update the user.
	fetchedUser, fetchedUserError := userUseCase.GetUserByEmail(ctx, email)
	if validator.IsErrorNotNil(fetchedUserError) {
		fetchedUserError = domainError.HandleError(fetchedUserError)
		return fetchedUserError
	}
	updatedUserPasswordError := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Email, "passwordResetToken", encodedTokenValue, "passwordResetAt", tokenExpirationTime)
	if validator.IsErrorNotNil(updatedUserPasswordError) {
		updatedUserPasswordError = domainError.HandleError(updatedUserPasswordError)
		return updatedUserPasswordError
	}

	templateName := config.ForgottenPasswordEmailTemplateName
	templatePath := config.ForgottenPasswordEmailTemplatePath
	emailData, prepareEmailDataError := PrepareEmailData(ctx, fetchedUser.Name, forgottenPasswordUrl, forgottenPasswordSubject, tokenValue, templateName, templatePath)
	if validator.IsErrorNotNil(prepareEmailDataError) {
		logging.Logger(prepareEmailDataError)
		prepareEmailDataError = domainError.HandleError(prepareEmailDataError)
		return prepareEmailDataError
	}

	sendEmailForgottenPasswordMessageError := userUseCase.userRepository.SendEmailForgottenPasswordMessage(ctx, fetchedUser, emailData)
	if validator.IsErrorNotNil(sendEmailForgottenPasswordMessageError) {
		sendEmailForgottenPasswordMessageError = domainError.HandleError(sendEmailForgottenPasswordMessageError)
		return sendEmailForgottenPasswordMessageError
	}
	return nil
}

func (userUseCase UserUseCase) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	updatedUser := userUseCase.userRepository.ResetUserPassword(ctx, firstKey, firstValue, secondKey, passwordKey, password)
	return updatedUser
}

func PrepareEmailData(ctx context.Context, userName string, emailUrl string, subject string,
	tokenValue string, templateName string, templatePath string) (userModel.EmailData, error) {
	loadConfig, loadConfigError := config.LoadConfig(".")
	if validator.IsErrorNotNil(loadConfigError) {
		var loadConfigInternalError domainError.InternalError
		loadConfigInternalError.Location = "User.Domain.UserUseCase.Registration.PrepareEmailData.LoadConfig"
		loadConfigInternalError.Reason = loadConfigError.Error()
		return userModel.EmailData{}, loadConfigInternalError
	}
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.EmailData{
		URL:          loadConfig.ClientOriginUrl + emailUrl + tokenValue,
		TemplateName: templateName,
		TemplatePath: templatePath,
		FirstName:    userFirstName,
		Subject:      subject,
	}
	return emailData, nil
}
