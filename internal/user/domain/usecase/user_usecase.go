package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/thanhpk/randstr"
	"github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
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

func (userUseCase *UserUseCase) GetAllUsers(ctx context.Context, page int, limit int) (*userModel.Users, error) {
	fetchedUsers, getAllUsers := userUseCase.userRepository.GetAllUsers(ctx, page, limit)
	return fetchedUsers, getAllUsers
}

func (userUseCase *UserUseCase) GetUserById(ctx context.Context, userID string) (*userModel.User, error) {
	fetchedUser, getUserByIdError := userUseCase.userRepository.GetUserById(ctx, userID)
	return fetchedUser, getUserByIdError
}

func (userUseCase *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	fetchedUser, getUserByEmailError := userUseCase.userRepository.GetUserByEmail(ctx, email)
	return fetchedUser, getUserByEmailError
}

func (userUseCase *UserUseCase) Register(ctx context.Context, userCreate *userModel.UserCreate) error {
	checkEmailDublicateError := userUseCase.userRepository.CheckEmailDublicate(ctx, userCreate.Email)
	if validator.IsErrorNotNil(checkEmailDublicateError) {
		checkEmailDublicateError = domainError.HandleError(checkEmailDublicateError)
		return checkEmailDublicateError
	}
	userCreateValidationErrors := UserCreateValidator(userCreate)
	if validator.IsErrorNotNil(userCreateValidationErrors) {
		userCreateValidationErrors = domainError.HandleError(userCreateValidationErrors)
		return userCreateValidationErrors
	}
	userCreate.Verified = true
	userCreate.Role = "user"
	createdUser := userUseCase.userRepository.Register(ctx, userCreate)
	if validator.IsErrorNotNil(createdUser.Error) {
		createdUser.Error = domainError.HandleError(createdUser.Error)
		return createdUser.Error
	}

	// Generate verification code.
	tokenValue := randstr.String(verificationCodeLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	_, userUpdateError := userUseCase.userRepository.UpdateNewRegisteredUserById(ctx, createdUser.Data.UserID, "verificationCode", encodedTokenValue)
	if validator.IsErrorNotNil(userUpdateError) {
		userUpdateError = domainError.HandleError(userUpdateError)
		return userUpdateError
	}
	templateName := config.UserConfirmationEmailTemplateName
	templatePath := config.UserConfirmationEmailTemplatePath

	emailData, prepareEmailDataError := PrepareEmailData(ctx, createdUser.Data.Name, emailConfirmationUrl, emailConfirmationSubject, tokenValue, templateName, templatePath)
	if validator.IsErrorNotNil(prepareEmailDataError) {
		logging.Logger(prepareEmailDataError)
		prepareEmailDataError = domainError.HandleError(prepareEmailDataError)
		return prepareEmailDataError
	}
	sendEmailVerificationMessageError := userUseCase.userRepository.SendEmailVerificationMessage(ctx, createdUser.Data, emailData)
	if validator.IsErrorNotNil(sendEmailVerificationMessageError) {
		sendEmailVerificationMessageError = domainError.HandleError(sendEmailVerificationMessageError)
		return sendEmailVerificationMessageError
	}
	return nil
}

func (userUseCase *UserUseCase) UpdateUserById(ctx context.Context, userID string, user *userModel.UserUpdate) (*userModel.User, error) {
	userUpdateValidationErrors := UserUpdateValidator(user)
	if validator.IsErrorNotNil(userUpdateValidationErrors) {
		userUpdateValidationErrors = domainError.HandleError(userUpdateValidationErrors)
		return nil, userUpdateValidationErrors
	}

	updatedUser, userUpdateError := userUseCase.userRepository.UpdateUserById(ctx, userID, user)
	if validator.IsErrorNotNil(userUpdateError) {
		userUpdateError = domainError.HandleError(userUpdateError)
		return nil, userUpdateError
	}
	return updatedUser, nil
}

func (userUseCase *UserUseCase) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCase.userRepository.DeleteUserById(ctx, userID)
	return deletedUser
}

func (userUseCase *UserUseCase) Login(ctx context.Context, userLogin *userModel.UserLogin) (string, error) {
	userLoginValidationErrors := UserLoginValidator(userLogin)
	if validator.IsErrorNotNil(userLoginValidationErrors) {
		userLoginValidationErrors = domainError.HandleError(userLoginValidationErrors)
		return "", userLoginValidationErrors
	}
	fetchedUser, getUserByEmailError := userUseCase.userRepository.GetUserByEmail(ctx, userLogin.Email)

	// Will return wrong email or password.
	if validator.IsErrorNotNil(getUserByEmailError) {
		return "", fmt.Errorf("invalid email or password")
	}

	// Verify password - we previously created this method.
	matchPasswords := domainUtility.VerifyPassword(fetchedUser.Password, userLogin.Password)
	if validator.IsErrorNotNil(matchPasswords) {
		return "", fmt.Errorf("invalid email or password")
	}
	return fetchedUser.UserID, nil
}

func (userUseCase *UserUseCase) UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*userModel.User, error) {
	updatedUser, updateNewRegisteredUserById := userUseCase.userRepository.UpdateNewRegisteredUserById(ctx, userID, key, value)
	return updatedUser, updateNewRegisteredUserById
}

func (userUseCase *UserUseCase) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
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

func (userUseCase *UserUseCase) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	updatedUser := userUseCase.userRepository.ResetUserPassword(ctx, firstKey, firstValue, secondKey, passwordKey, password)
	return updatedUser
}

func PrepareEmailData(ctx context.Context, userName string, emailUrl string, subject string,
	tokenValue string, templateName string, templatePath string) (*userModel.EmailData, error) {
	loadConfig, loadConfigError := config.LoadConfig(".")
	if validator.IsErrorNotNil(loadConfigError) {
		var loadConfigInternalError *domainError.InternalError = new(domainError.InternalError)
		loadConfigInternalError.Location = "User.Domain.UserUseCase.Registration.PrepareEmailData.LoadConfig"
		loadConfigInternalError.Reason = loadConfigError.Error()
		return nil, loadConfigInternalError
	}
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := &userModel.EmailData{
		URL:          loadConfig.ClientOriginUrl + emailUrl + tokenValue,
		TemplateName: templateName,
		TemplatePath: templatePath,
		FirstName:    userFirstName,
		Subject:      subject,
	}
	return emailData, nil
}
