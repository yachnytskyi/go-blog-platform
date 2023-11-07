package usecase

import (
	"context"
	"time"

	"github.com/thanhpk/randstr"
	config "github.com/yachnytskyi/golang-mongo-grpc/config"
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
	userRole                        = "user"
)

type UserUseCase struct {
	userRepository user.UserRepository
}

func NewUserUseCase(userRepository user.UserRepository) user.UserUseCase {
	return UserUseCase{userRepository: userRepository}
}

// GetAllUsers retrieves a list of users based on the provided pagination parameters.
// It returns a Result containing user data on success or an error on failure.
func (userUseCase UserUseCase) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	fetchedUsers := userUseCase.userRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsErrorNotNil(fetchedUsers.Error) {
		return commonModel.NewResultOnFailure[userModel.Users](domainError.HandleError(fetchedUsers.Error))
	}
	return fetchedUsers
}

// GetUserById retrieves a user by their ID using the provided context and user ID.
// It returns a Result containing user data on success or an error on failure.
func (userUseCase UserUseCase) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	fetchedUser := userUseCase.userRepository.GetUserById(ctx, userID)
	if validator.IsErrorNotNil(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}
	return fetchedUser
}

func (userUseCase UserUseCase) GetUserByEmail(ctx context.Context, email string) (userModel.User, error) {
	validateEmailError := validateEmail(email, emailRegex)
	if validator.IsValueNotNil(validateEmailError) {
		validateEmailError := domainError.HandleError(validateEmailError)
		return userModel.User{}, validateEmailError
	}
	fetchedUser, getUserByEmailError := userUseCase.userRepository.GetUserByEmail(ctx, email)
	return fetchedUser, getUserByEmailError
}

// Register is a method in the UserUseCase that handles the registration of a new user.
// It returns a Result containing created user data on success or an error on failure.
func (userUseCase UserUseCase) Register(ctx context.Context, userCreateData userModel.UserCreate) commonModel.Result[userModel.User] {
	// Step 1: Validate the user creation data.
	userCreate := validateUserCreate(userCreateData)
	if validator.IsErrorNotNil(userCreate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userCreate.Error))
	}

	// Step 2: Check for duplicate email.
	checkEmailDuplicateError := userUseCase.userRepository.CheckEmailDuplicate(ctx, userCreate.Data.Email)
	if validator.IsErrorNotNil(checkEmailDuplicateError) {
		checkEmailDuplicateError = domainError.HandleError(checkEmailDuplicateError)
		return commonModel.NewResultOnFailure[userModel.User](checkEmailDuplicateError)
	}

	// Step 3: Generate a verification token and set user properties.
	tokenValue := randstr.String(verificationCodeLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = encodedTokenValue

	// Step 4: Register the user.
	createdUser := userUseCase.userRepository.Register(ctx, userCreate.Data)
	if validator.IsErrorNotNil(createdUser.Error) {
		createdUser.Error = domainError.HandleError(createdUser.Error)
		return createdUser
	}

	// Step 5: Prepare email data for user registration.
	emailData := prepareEmailDataForUserRegister(ctx, createdUser.Data.Name, tokenValue)
	if validator.IsErrorNotNil(emailData.Error) {
		logging.Logger(emailData.Error)
		emailData.Error = domainError.HandleError(emailData.Error)
		return commonModel.NewResultOnFailure[userModel.User](emailData.Error)
	}

	// Step 6: Send the email verification message and return the created user.
	sendEmailVerificationMessageError := userUseCase.userRepository.SendEmailVerificationMessage(ctx, createdUser.Data, emailData.Data)
	if validator.IsErrorNotNil(sendEmailVerificationMessageError) {
		sendEmailVerificationMessageError = domainError.HandleError(sendEmailVerificationMessageError)
		return commonModel.NewResultOnFailure[userModel.User](sendEmailVerificationMessageError)
	}
	return createdUser
}

func (userUseCase UserUseCase) UpdateUserById(ctx context.Context, userUpdateData userModel.UserUpdate) (userModel.User, error) {
	userUpdate := validateUserUpdate(userUpdateData)
	if validator.IsErrorNotNil(userUpdate.Error) {
		validationErrors := domainError.HandleError(userUpdate.Error)
		return userModel.User{}, validationErrors
	}

	updatedUser, userUpdateError := userUseCase.userRepository.UpdateUserById(ctx, userUpdate.Data)
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
		return "", getUserByEmailError
	}
	arePasswordsNotEqualError := arePasswordsNotEqual(fetchedUser.Password, userLoginData.Password)
	if validator.IsValueNotNil(arePasswordsNotEqualError) {
		arePasswordsNotEqualError.Notification = invalidEmailOrPassword
		return "", arePasswordsNotEqualError
	}
	return fetchedUser.UserID, nil
}

func (userUseCase UserUseCase) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {

	validateEmailError := validateEmail(email, emailRegex)
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

	emailData := prepareEmailDataForUserUpdate(ctx, fetchedUser.Name, tokenValue)
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

func prepareEmailDataForUserRegister(ctx context.Context, userName string, tokenValue string) commonModel.Result[userModel.EmailData] {
	applicationConfig := config.AppConfig
	subject := emailConfirmationSubject
	url := emailConfirmationUrl
	templateName := applicationConfig.Email.UserConfirmationTemplateName
	templatePath := applicationConfig.Email.UserConfirmationTemplatePath
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(applicationConfig.Email.ClientOriginUrl+url+tokenValue, templateName, templatePath, userFirstName, subject)
	return commonModel.NewResultOnSuccess[userModel.EmailData](emailData)
}

func prepareEmailDataForUserUpdate(ctx context.Context, userName string, tokenValue string) commonModel.Result[userModel.EmailData] {
	applicationConfig := config.AppConfig
	subject := forgottenPasswordSubject
	url := forgottenPasswordUrl
	templateName := applicationConfig.Email.ForgottenPasswordTemplateName
	templatePath := applicationConfig.Email.ForgottenPasswordTemplatePath
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(applicationConfig.Email.ClientOriginUrl+url+tokenValue, templateName, templatePath, userFirstName, subject)
	return commonModel.NewResultOnSuccess[userModel.EmailData](emailData)
}
