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
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/error/domain"
	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	verificationCodeLength int = 20
	resetTokenLength       int = 20
	userRole                   = "user"
)

type UserUseCase struct {
	userRepository user.UserRepository
}

func NewUserUseCase(userRepository user.UserRepository) user.UserUseCase {
	return UserUseCase{userRepository: userRepository}
}

// GetAllUsers retrieves a list of users based on the provided pagination parameters.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	// Fetch the users.
	fetchedUsers := userUseCase.userRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		return commonModel.NewResultOnFailure[userModel.Users](domainError.HandleError(fetchedUsers.Error))
	}
	return fetchedUsers
}

// GetUserById retrieves a user by their ID using the user ID.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	// Check context timeout.
	contextError := commonUtility.HandleWithContextError("internal.user.domain.usecase.GetUserById", ctx)
	if validator.IsError(contextError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(contextError))
	}

	// Fetch the user.
	fetchedUser := userUseCase.userRepository.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}
	return fetchedUser
}

// GetUserByEmail retrieves a user by their ID using the provided user email.
// It performs email format validation and fetches the user from the repository.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	// Validate the email.
	validateEmailError := validateEmail(email, emailRegex)
	if validator.IsError(validateEmailError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(validateEmailError))
	}

	// Fetch the user.
	fetchedUser := userUseCase.userRepository.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}
	return fetchedUser
}

// Register registers a new user based on the provided data, generates a verification token,
// and sends an email verification message. The result is wrapped in a commonModel.Result
// containing either the user or an error.
func (userUseCase UserUseCase) Register(ctx context.Context, userCreateData userModel.UserCreate) commonModel.Result[userModel.User] {
	// Validate the user creation data.
	userCreate := validateUserCreate(userCreateData)
	if validator.IsError(userCreate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userCreate.Error))
	}

	// Check for duplicate email.
	checkEmailDuplicateError := userUseCase.userRepository.CheckEmailDuplicate(ctx, userCreate.Data.Email)
	if validator.IsError(checkEmailDuplicateError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(checkEmailDuplicateError))
	}

	// Generate a verification token and set user properties.
	tokenValue := randstr.String(verificationCodeLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = encodedTokenValue
	currentTime := time.Now()
	userCreate.Data.CreatedAt = currentTime
	userCreate.Data.UpdatedAt = currentTime

	// Register the user.
	createdUser := userUseCase.userRepository.Register(ctx, userCreate.Data)
	if validator.IsError(createdUser.Error) {
		createdUser.Error = domainError.HandleError(createdUser.Error)
		return createdUser
	}

	// Prepare email data for user registration.
	// Send the email verification message and return the created user.
	emailData := prepareEmailDataForRegistration(ctx, createdUser.Data.Name, tokenValue)
	sendEmailVerificationMessageError := userUseCase.userRepository.SendEmailVerificationMessage(ctx, createdUser.Data, emailData)
	if validator.IsError(sendEmailVerificationMessageError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(sendEmailVerificationMessageError))
	}
	return createdUser
}

// UpdateUserById updates a user's information based on the provided data.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) UpdateCurrentUser(ctx context.Context, userUpdateData userModel.UserUpdate) commonModel.Result[userModel.User] {
	// Validate the user update data.
	userUpdate := validateUserUpdate(userUpdateData)
	if validator.IsError(userUpdate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userUpdate.Error))
	}

	// Set the update timestamp to the current time.
	userUpdate.Data.UpdatedAt = time.Now()

	// Update the user.
	updatedUser := userUseCase.userRepository.UpdateCurrentUser(ctx, userUpdate.Data)
	if validator.IsError(updatedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(updatedUser.Error))
	}
	return updatedUser
}

// DeleteUserById deletes a user based on the provided user ID.
func (userUseCase UserUseCase) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCase.userRepository.DeleteUserById(ctx, userID)
	if validator.IsError(deletedUser) {
		return domainError.HandleError(deletedUser)
	}

	// User deletion was successful. Return nil to indicate no error.
	return nil
}

// Login performs the user authentication process.
// It takes a userLoginData object, validates it, fetches the user by email,
// checks if the provided password matches the stored password, and generates
// access and refresh tokens upon successful authentication.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCase UserUseCase) Login(ctx context.Context, userLoginData userModel.UserLogin) commonModel.Result[userModel.UserToken] {
	// Validate the user login data.
	userLogin := validateUserLogin(userLoginData)
	if validator.IsError(userLogin.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userLogin.Error))
	}

	// Fetch the user by email from the repository.
	fetchedUser := userUseCase.userRepository.GetUserByEmail(ctx, userLogin.Data.Email)

	// Check if the provided password matches the stored password.
	arePasswordsNotEqualError := arePasswordsNotEqual(fetchedUser.Data.Password, userLoginData.Password)
	if validator.IsError(arePasswordsNotEqualError) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(arePasswordsNotEqualError))
	}

	// Generate access and refresh tokens.
	userToken := generateToken(ctx, fetchedUser.Data.UserID)

	// Return the result containing userToken information.
	return userToken
}

func (userUseCase UserUseCase) RefreshAccessToken(ctx context.Context, user userModel.User) commonModel.Result[userModel.UserToken] {
	// Generate access and refresh tokens.
	userToken := generateToken(ctx, user.UserID)

	// Return the result containing userToken information.
	return userToken
}

func (userUseCase UserUseCase) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {

	validateEmailError := validateEmail(email, emailRegex)
	if validator.IsError(validateEmailError) {
		return domainError.HandleError(validateEmailError)
	}
	updatedUserError := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, email, firstKey, firstValue, secondKey, secondValue)
	if validator.IsError(updatedUserError) {
		updatedUserError = domainError.HandleError(updatedUserError)
		return updatedUserError
	}

	// Generate verification code.
	tokenValue := randstr.String(resetTokenLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	tokenExpirationTime := time.Now().Add(time.Minute * 15)

	// Update the user.
	fetchedUser := userUseCase.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		fetchedUserError := domainError.HandleError(fetchedUser.Error)
		return fetchedUserError
	}
	updatedUserPasswordError := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Data.Email, "passwordResetToken", encodedTokenValue, "passwordResetAt", tokenExpirationTime)
	if validator.IsError(updatedUserPasswordError) {
		updatedUserPasswordError = domainError.HandleError(updatedUserPasswordError)
		return updatedUserPasswordError
	}

	emailData := prepareEmailDataForUpdatePasswordResetToken(ctx, fetchedUser.Data.Name, tokenValue)
	sendEmailForgottenPasswordMessageError := userUseCase.userRepository.SendEmailForgottenPasswordMessage(ctx, fetchedUser.Data, emailData)
	if validator.IsError(sendEmailForgottenPasswordMessageError) {
		sendEmailForgottenPasswordMessageError = domainError.HandleError(sendEmailForgottenPasswordMessageError)
		return sendEmailForgottenPasswordMessageError
	}
	return nil
}

func (userUseCase UserUseCase) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	updatedUser := userUseCase.userRepository.ResetUserPassword(ctx, firstKey, firstValue, secondKey, passwordKey, password)
	return updatedUser
}

// prepareEmailData is a helper function to create an EmailData model for sending an email.
// It takes the context, user name, token value, email subject, URL, template name, and template path as input.
// It constructs an EmailData model and returns it in a Result.
func prepareEmailData(ctx context.Context, userName, tokenValue, subject, url, templateName, templatePath string) userModel.EmailData {
	emailConfig := config.AppConfig.Email
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(emailConfig.ClientOriginUrl+url+tokenValue, templateName, templatePath, userFirstName, subject)
	return emailData
}

// prepareEmailDataForRegister is a helper function to prepare an EmailData model specifically for user registration.
// It takes the context, user name, and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
func prepareEmailDataForRegistration(ctx context.Context, userName, tokenValue string) userModel.EmailData {
	emailConfig := config.AppConfig.Email
	return prepareEmailData(ctx, userName, tokenValue, constants.EmailConfirmationSubject, constants.EmailConfirmationUrl,
		emailConfig.UserConfirmationTemplateName, emailConfig.UserConfirmationTemplatePath)
}

// prepareEmailDataForUserUpdate is a helper function to prepare an EmailData model specifically for updating user information.
// It takes the context, user name, and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
func prepareEmailDataForUpdatePasswordResetToken(ctx context.Context, userName, tokenValue string) userModel.EmailData {
	emailConfig := config.AppConfig.Email
	return prepareEmailData(ctx, userName, tokenValue, constants.ForgottenPasswordSubject, constants.ForgottenPasswordUrl,
		emailConfig.ForgottenPasswordTemplateName, emailConfig.ForgottenPasswordTemplatePath)
}

// generateToken generates access and refresh tokens for a user.
// It takes the user ID as input and uses the application configuration
// to create both access and refresh tokens using domainUtility.GenerateJWTToken.
// If there are errors during token generation, it returns a failure result
// with the corresponding error. Otherwise, it returns a success result
// containing the generated access and refresh tokens.
func generateToken(ctx context.Context, userID string) commonModel.Result[userModel.UserToken] {
	// Create a userToken struct to store the generated tokens.
	var userToken userModel.UserToken

	// Retrieve application configuration.
	accessTokenConfig := config.AppConfig.AccessToken
	refreshTokenConfig := config.AppConfig.RefreshToken

	// Generate the access token.
	accessToken, accessTokenGenerationError := domainUtility.GenerateJWTToken(ctx, accessTokenConfig.ExpiredIn, userID, accessTokenConfig.PrivateKey)
	if validator.IsError(accessTokenGenerationError) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(accessTokenGenerationError))
	}

	// Generate the refresh token.
	refreshToken, refreshTokenGenerationError := domainUtility.GenerateJWTToken(ctx, refreshTokenConfig.ExpiredIn, userID, refreshTokenConfig.PrivateKey)
	if validator.IsError(refreshTokenGenerationError) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(refreshTokenGenerationError))
	}

	// Update the userToken struct with the generated tokens.
	userToken.AccessToken = accessToken
	userToken.RefreshToken = refreshToken

	// Return a success result with the userToken struct.
	return commonModel.NewResultOnSuccess[userModel.UserToken](userToken)
}
