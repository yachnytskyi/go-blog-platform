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

// GetAllUsers retrieves a list of users based on the provided pagination parameters.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCaseV1 UserUseCaseV1) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	// Fetch the users.
	fetchedUsers := userUseCaseV1.userRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		return commonModel.NewResultOnFailure[userModel.Users](domainError.HandleError(fetchedUsers.Error))
	}

	return fetchedUsers
}

// GetUserById retrieves a user by their ID using the user ID.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCaseV1 UserUseCaseV1) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	// Fetch the user.
	fetchedUser := userUseCaseV1.userRepository.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

// GetUserByEmail retrieves a user by their ID using the provided user email.
// It performs email format validation and fetches the user from the repository.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCaseV1 UserUseCaseV1) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	// Validate the email.
	validateEmailError := isEmailValid(email)
	if validator.IsError(validateEmailError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(validateEmailError))
	}

	// Fetch the user.
	fetchedUser := userUseCaseV1.userRepository.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

// Register registers a new user based on the provided data, generates a verification token,
// and sends an email verification message. The result is wrapped in a commonModel.Result
// containing either the user or an error.
func (userUseCaseV1 UserUseCaseV1) Register(ctx context.Context, userCreateData userModel.UserCreate) commonModel.Result[userModel.User] {
	// Validate the user creation data.
	userCreate := validateUserCreate(userCreateData)
	if validator.IsError(userCreate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userCreate.Error))
	}

	// Check for duplicate email.
	checkEmailDuplicateError := userUseCaseV1.userRepository.CheckEmailDuplicate(ctx, userCreate.Data.Email)
	if validator.IsError(checkEmailDuplicateError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(checkEmailDuplicateError))
	}

	// Generate a verification token and set user properties.
	token := randstr.String(verificationCodeLength)
	token = commonUtility.Encode(token)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = token
	currentTime := time.Now()
	userCreate.Data.CreatedAt = currentTime
	userCreate.Data.UpdatedAt = currentTime

	// Register the user.
	createdUser := userUseCaseV1.userRepository.Register(ctx, userCreate.Data)
	if validator.IsError(createdUser.Error) {
		createdUser.Error = domainError.HandleError(createdUser.Error)
		return createdUser
	}

	// Prepare email data for user registration.
	// Send the email verification message and return the created user.
	emailData := prepareEmailDataForRegistration(createdUser.Data.Name, token)
	sendEmailError := userUseCaseV1.userRepository.SendEmail(createdUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(sendEmailError))
	}

	return createdUser
}

// UpdateUserById updates a user's information based on the provided data.
// The result is wrapped in a commonModel.Result containing either the user or an error.
func (userUseCaseV1 UserUseCaseV1) UpdateCurrentUser(ctx context.Context, userUpdateData userModel.UserUpdate) commonModel.Result[userModel.User] {
	// Validate the user update data.
	userUpdate := validateUserUpdate(userUpdateData)
	if validator.IsError(userUpdate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userUpdate.Error))
	}

	// Set the update timestamp to the current time.
	userUpdate.Data.UpdatedAt = time.Now()

	// Update the user.
	updatedUser := userUseCaseV1.userRepository.UpdateCurrentUser(ctx, userUpdate.Data)
	if validator.IsError(updatedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(updatedUser.Error))
	}

	return updatedUser
}

// DeleteUserById deletes a user based on the provided user ID.
func (userUseCaseV1 UserUseCaseV1) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCaseV1.userRepository.DeleteUserById(ctx, userID)
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
func (userUseCaseV1 UserUseCaseV1) Login(ctx context.Context, userLoginData userModel.UserLogin) commonModel.Result[userModel.UserToken] {
	// Validate the user login data.
	userLogin := validateUserLogin(userLoginData)
	if validator.IsError(userLogin.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userLogin.Error))
	}

	// Fetch the user by email from the repository.
	fetchedUser := userUseCaseV1.userRepository.GetUserByEmail(ctx, userLogin.Data.Email)

	// Check if the provided password matches the stored password.
	arePasswordsNotEqualError := arePasswordsNotEqual(fetchedUser.Data.Password, userLoginData.Password)
	if validator.IsError(arePasswordsNotEqualError) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(arePasswordsNotEqualError))
	}

	// Generate the UserTokenPayload.
	userTokenPayload := domainModel.NewUserTokenPayload(fetchedUser.Data.ID, fetchedUser.Data.Role)

	// Generate access and refresh tokens.
	userToken := generateToken(userTokenPayload)
	if validator.IsError(userLogin.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userToken.Error))
	}

	// Return the result containing userToken information.
	return userToken
}

// RefreshAccessToken generates a new access token using the provided user information.
// It creates a UserTokenPayload and uses it to generate new access and refresh tokens.
// The result is wrapped in a commonModel.Result containing either the tokens or an error.
func (userUseCaseV1 UserUseCaseV1) RefreshAccessToken(ctx context.Context, user userModel.User) commonModel.Result[userModel.UserToken] {
	// Generate the UserTokenPayload using the user ID and role.
	userTokenPayload := domainModel.NewUserTokenPayload(user.ID, user.Role)

	// Generate new access and refresh tokens.
	userToken := generateToken(userTokenPayload)

	// Return the result containing the updated token information.
	return userToken
}

// ForgottenPassword handles the process of initiating a password reset for a user.
// It validates the email, generates a reset token, updates the user record with the token and expiration,
// and sends an email with the reset instructions.
func (userUseCaseV1 UserUseCaseV1) ForgottenPassword(ctx context.Context, userForgottenPassword userModel.UserForgottenPassword) error {
	// Validate the user forgotten password data.
	validateEmailError := validateUserForgottenPassword(userForgottenPassword)
	if validator.IsError(validateEmailError.Error) {
		return domainError.HandleError(validateEmailError.Error)
	}

	// Fetch the user by email.
	fetchedUser := userUseCaseV1.GetUserByEmail(ctx, userForgottenPassword.Email)
	if validator.IsError(fetchedUser.Error) {
		fetchedUserError := domainError.HandleError(fetchedUser.Error)
		return fetchedUserError
	}

	// Generate a reset token and set the expiration time.
	tokenValue := randstr.String(resetTokenLength)
	tokenValue = commonUtility.Encode(tokenValue)
	tokenExpirationTime := time.Now().Add(constants.PasswordResetTokenExpirationTime)

	// Set the reset token and expiration in the user model.
	userForgottenPassword.ResetToken = tokenValue
	userForgottenPassword.ResetExpiry = tokenExpirationTime

	// Update the user record with the user forgotten password data.
	updatedUserPasswordError := userUseCaseV1.userRepository.ForgottenPassword(ctx, userForgottenPassword)
	if validator.IsError(updatedUserPasswordError) {
		updatedUserPasswordError = domainError.HandleError(updatedUserPasswordError)
		return updatedUserPasswordError
	}

	// Prepare and send the reset email.
	emailData := prepareEmailDataForUpdatePasswordResetToken(fetchedUser.Data.Name, tokenValue)
	sendEmailError := userUseCaseV1.userRepository.SendEmail(fetchedUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return domainError.HandleError(sendEmailError)
	}

	// Return nil to indicate success.
	return nil
}

func (userUseCaseV1 UserUseCaseV1) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	updatedUser := userUseCaseV1.userRepository.ResetUserPassword(ctx, firstKey, firstValue, secondKey, passwordKey, password)
	return updatedUser
}

// prepareEmailData is a helper function to create an EmailData model for sending an email.
// It takes the context, user name, token value, email subject, URL, template name, and template path as input.
// It constructs an EmailData model and returns it in a Result.
func prepareEmailData(userName, tokenValue, subject, url, templateName, templatePath string) userModel.EmailData {
	emailConfig := config.GetEmailConfig()
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(emailConfig.ClientOriginUrl+url+tokenValue, templateName, templatePath, userFirstName, subject)
	return emailData
}

// prepareEmailDataForRegister is a helper function to prepare an EmailData model specifically for user registration.
// It takes the context, user name, and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
func prepareEmailDataForRegistration(userName, tokenValue string) userModel.EmailData {
	emailConfig := config.GetEmailConfig()
	return prepareEmailData(userName, tokenValue, constants.EmailConfirmationSubject, constants.EmailConfirmationUrl,
		emailConfig.UserConfirmationTemplateName, emailConfig.UserConfirmationTemplatePath)
}

// prepareEmailDataForUserUpdate is a helper function to prepare an EmailData model specifically for updating user information.
// It takes the context, user name, and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
func prepareEmailDataForUpdatePasswordResetToken(userName, tokenValue string) userModel.EmailData {
	emailConfig := config.GetEmailConfig()
	return prepareEmailData(userName, tokenValue, constants.ForgottenPasswordSubject, constants.ForgottenPasswordUrl,
		emailConfig.ForgottenPasswordTemplateName, emailConfig.ForgottenPasswordTemplatePath)
}

// generateToken generates access and refresh tokens for a user.
// It takes the user ID as input and uses the application configuration
// to create both access and refresh tokens using domainUtility.GenerateJWTToken.
// If there are errors during token generation, it returns a failure result
// with the corresponding error. Otherwise, it returns a success result
// containing the generated access and refresh tokens.
func generateToken(userTokenPayload domainModel.UserTokenPayload) commonModel.Result[userModel.UserToken] {
	// Create a userToken struct to store the generated tokens.
	var userToken userModel.UserToken

	// Retrieve application configuration.
	accessTokenConfig := config.GetAccessConfig()
	refreshTokenConfig := config.GetRefreshConfig()

	// Generate the access token.
	accessToken := domainUtility.GenerateJWTToken(location+".generateToken.accessToken", accessTokenConfig.PrivateKey, accessTokenConfig.ExpiredIn, userTokenPayload)
	if validator.IsError(accessToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](accessToken.Error)
	}

	// Generate the refresh token.
	refreshToken := domainUtility.GenerateJWTToken(location+".generateToken.refreshToken", refreshTokenConfig.PrivateKey, refreshTokenConfig.ExpiredIn, userTokenPayload)
	if validator.IsError(accessToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](refreshToken.Error)
	}

	// Update the userToken struct with the generated tokens.
	userToken.AccessToken = accessToken.Data
	userToken.RefreshToken = refreshToken.Data

	// Return a success result with the userToken struct.
	return commonModel.NewResultOnSuccess[userModel.UserToken](userToken)
}
