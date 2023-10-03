package usecase

import (
	"context"
	"time"

	"github.com/thanhpk/randstr"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
	"github.com/yachnytskyi/golang-mongo-grpc/config/constant"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
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
	userRepository user.UserRepository
}

func NewUserUseCase(userRepository user.UserRepository) user.UserUseCase {
	return &UserUseCase{userRepository: userRepository}
}

func (userUseCase UserUseCase) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	fetchedUsers := userUseCase.userRepository.GetAllUsers(ctx, paginationQuery)
	return fetchedUsers
}

func (userUseCase UserUseCase) GetUserById(ctx context.Context, userID string) (userModel.User, error) {
	fetchedUser, getUserByIdError := userUseCase.userRepository.GetUserById(ctx, userID)
	return fetchedUser, getUserByIdError
}

func (userUseCase UserUseCase) GetUserByEmail(ctx context.Context, email string) (userModel.User, error) {
	validateEmailError := validateEmail(email, constant.FieldRequired, emailRegex)
	if validator.IsValueNotNil(validateEmailError) {
		validateEmailError := domainError.HandleError(validateEmailError)
		return userModel.User{}, validateEmailError
	}
	fetchedUser, getUserByEmailError := userUseCase.userRepository.GetUserByEmail(ctx, email)
	return fetchedUser, getUserByEmailError
}

func (userUseCase UserUseCase) Register(ctx context.Context, userCreateData userModel.UserCreate) commonModel.Result[userModel.User] {
	userCreate := validateUserCreate(userCreateData)
	if validator.IsErrorNotNil(userCreate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userCreate.Error))
	}
	checkEmailDublicateError := userUseCase.userRepository.CheckEmailDublicate(ctx, userCreate.Data.Email)
	if validator.IsErrorNotNil(checkEmailDublicateError) {
		checkEmailDublicateError = domainError.HandleError(checkEmailDublicateError)
		return commonModel.NewResultOnFailure[userModel.User](checkEmailDublicateError)
	}
	tokenValue := randstr.String(verificationCodeLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	userCreate.Data.Role = "user"
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = encodedTokenValue

	createdUser := userUseCase.userRepository.Register(ctx, userCreate.Data)
	if validator.IsErrorNotNil(createdUser.Error) {
		createdUser.Error = domainError.HandleError(createdUser.Error)
		return createdUser
	}
	applicationConfig := config.AppConfig
	templateName := applicationConfig.Email.UserConfirmationTemplateName
	templatePath := applicationConfig.Email.UserConfirmationTemplatePath

	emailData := PrepareEmailData(ctx, createdUser.Data.Name, emailConfirmationUrl, emailConfirmationSubject, tokenValue, templateName, templatePath)
	if validator.IsErrorNotNil(emailData.Error) {
		logging.Logger(emailData.Error)
		emailData.Error = domainError.HandleError(emailData.Error)
		return commonModel.NewResultOnFailure[userModel.User](emailData.Error)
	}
	sendEmailVerificationMessageError := userUseCase.userRepository.SendEmailVerificationMessage(ctx, createdUser.Data, emailData.Data)
	if validator.IsErrorNotNil(sendEmailVerificationMessageError) {
		sendEmailVerificationMessageError = domainError.HandleError(sendEmailVerificationMessageError)
		return commonModel.NewResultOnFailure[userModel.User](sendEmailVerificationMessageError)
	}
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

func (userUseCase UserUseCase) DeleteUser(ctx context.Context, userID string) error {
	deletedUser := userUseCase.userRepository.DeleteUser(ctx, userID)
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
	if arePasswordsNotEqual(fetchedUser.Password, userLoginData.Password) {
		return "", domainError.HandleError(domainError.NewErrorMessage(invalidEmailOrPassword))
	}
	return fetchedUser.UserID, nil
}

func (userUseCase UserUseCase) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {

	validateEmailError := validateEmail(email, EmailField, emailRegex)
	if validator.IsValueNotNil(validateEmailError) {
		validateEmailError := domainError.HandleError(validateEmailError)
		return validateEmailError
	}
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

	applicationConfig := config.AppConfig
	templateName := applicationConfig.Email.ForgottenPasswordTemplateName
	templatePath := applicationConfig.Email.ForgottenPasswordTemplatePath
	emailData := PrepareEmailData(ctx, fetchedUser.Name, forgottenPasswordUrl, forgottenPasswordSubject, tokenValue, templateName, templatePath)
	if validator.IsErrorNotNil(emailData.Error) {
		logging.Logger(emailData.Error)
		emailData.Error = domainError.HandleError(emailData.Error)
		return emailData.Error
	}

	sendEmailForgottenPasswordMessageError := userUseCase.userRepository.SendEmailForgottenPasswordMessage(ctx, fetchedUser, emailData.Data)
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

func PrepareEmailData(ctx context.Context, userName string, url string, subject string,
	tokenValue string, templateName string, templatePath string) commonModel.Result[userModel.EmailData] {
	applicationConfig := config.AppConfig
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(applicationConfig.Email.ClientOriginUrl+url+tokenValue, templateName, templatePath, userFirstName, subject)
	return commonModel.NewResultOnSuccess[userModel.EmailData](emailData)

}
