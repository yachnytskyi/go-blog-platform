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
	logging "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/logging"
	validator "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/validator"
)

const (
	location                   = "internal.user.domain.usecase."
	verificationCodeLength int = 20
	resetTokenLength       int = 20
	userRole                   = "user"
)

// UserUseCaseV1 implements the UserUseCase interface and provides methods to handle user-related operations.
type UserUseCaseV1 struct {
	userRepository user.UserRepository
}

// NewUserUseCaseV1 creates a new instance of UserUseCaseV1 with the given user repository.
//
// Parameters:
// - userRepository (user.UserRepository): The user repository instance that handles data access operations.
//
// Returns:
// - UserUseCaseV1: The initialized UserUseCaseV1 instance.
func NewUserUseCaseV1(userRepository user.UserRepository) user.UserUseCase {
	return UserUseCaseV1{userRepository: userRepository}
}

// GetAllUsers retrieves a list of users based on the provided pagination parameters.
// It performs the following steps:
// 1. Calls the user repository to fetch the list of users based on the pagination query.
// 2. Checks for errors during user fetching and wraps the error in a commonModel.Result if any.
// 3. Returns the fetched users wrapped in a commonModel.Result.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - paginationQuery (commonModel.PaginationQuery): The query parameters for pagination.
//
// Returns:
// - commonModel.Result[userModel.Users]: The result containing either the list of users or an error.
func (userUseCaseV1 UserUseCaseV1) GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users] {
	// Call the user repository to fetch the users.
	fetchedUsers := userUseCaseV1.userRepository.GetAllUsers(ctx, paginationQuery)
	if validator.IsError(fetchedUsers.Error) {
		return commonModel.NewResultOnFailure[userModel.Users](domainError.HandleError(fetchedUsers.Error))
	}

	// Return the fetched users wrapped in a commonModel.Result.
	return fetchedUsers
}

// GetUserById retrieves a user by their ID.
// It performs the following steps:
// 1. Calls the user repository to fetch the user based on the provided user ID.
// 2. Checks for errors during user fetching and wraps the error in a commonModel.Result if any.
// 3. Returns the fetched user wrapped in a commonModel.Result.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userID (string): The unique identifier of the user to be fetched.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
func (userUseCaseV1 UserUseCaseV1) GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User] {
	// Call the user repository to fetch the user.
	fetchedUser := userUseCaseV1.userRepository.GetUserById(ctx, userID)
	if validator.IsError(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

// GetUserByEmail retrieves a user by their email.
// It performs the following steps:
// 1. Validates the provided email format.
// 2. Calls the user repository to fetch the user based on the email.
// 3. Checks for errors during user fetching and wraps the error in a commonModel.Result if any.
// 4. Returns the fetched user wrapped in a commonModel.Result.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - email (string): The email of the user to be fetched.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
func (userUseCaseV1 UserUseCaseV1) GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User] {
	// Validate the email.
	validateEmailError := isEmailValid(email)
	if validator.IsError(validateEmailError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(validateEmailError))
	}

	// Call the user repository to fetch the user by email.
	fetchedUser := userUseCaseV1.userRepository.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(fetchedUser.Error))
	}

	return fetchedUser
}

// Register registers a new user based on the provided data, generates a verification token, and sends an email verification message.
// It performs the following steps:
// 1. Validates the provided user creation data.
// 2. Checks for a duplicate email in the repository.
// 3. Generates a verification token and sets user properties.
// 4. Registers the user in the repository.
// 5. Prepares email data and sends the email verification message.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userCreateData (userModel.UserCreate): The data for creating the user.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
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
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(createdUser.Error))
	}

	// Prepare email data for user registration and send the email verification message.
	emailData := prepareEmailDataForRegistration(createdUser.Data.Name, token)
	sendEmailError := userUseCaseV1.userRepository.SendEmail(createdUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(sendEmailError))
	}

	return createdUser
}

// UpdateCurrentUser updates a user's information based on the provided data.
// It performs the following steps:
// 1. Validates the provided user update data.
// 2. Sets the update timestamp to the current time.
// 3. Calls the user repository to update the user.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userUpdateData (userModel.UserUpdate): The data for updating the user.
//
// Returns:
// - commonModel.Result[userModel.User]: The result containing either the user data or an error.
func (userUseCaseV1 UserUseCaseV1) UpdateCurrentUser(ctx context.Context, userUpdateData userModel.UserUpdate) commonModel.Result[userModel.User] {
	// Validate the user update data.
	userUpdate := validateUserUpdate(userUpdateData)
	if validator.IsError(userUpdate.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(userUpdate.Error))
	}

	// Set the update timestamp to the current time.
	userUpdate.Data.UpdatedAt = time.Now()

	// Call the user repository to update the user.
	updatedUser := userUseCaseV1.userRepository.UpdateCurrentUser(ctx, userUpdate.Data)
	if validator.IsError(updatedUser.Error) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(updatedUser.Error))
	}

	return updatedUser
}

// DeleteUserById deletes a user based on the provided user ID.
// It performs the following steps:
// 1. Calls the user repository to delete the user based on the provided user ID.
// 2. Checks for errors during user deletion and handles the error if any.
// 3. Returns nil if the deletion was successful, indicating no error.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userID (string): The unique identifier of the user to be deleted.
//
// Returns:
// - error: An error if any occurred during user deletion, otherwise nil.
func (userUseCaseV1 UserUseCaseV1) DeleteUserById(ctx context.Context, userID string) error {
	// Call the user repository to delete the user.
	deletedUser := userUseCaseV1.userRepository.DeleteUserById(ctx, userID)
	if validator.IsError(deletedUser) {
		return domainError.HandleError(deletedUser)
	}

	// User deletion was successful. Return nil to indicate no error.
	return nil
}

// Login performs the user authentication process.
// It performs the following steps:
// 1. Validates the provided user login data.
// 2. Calls the user repository to fetch the user by email.
// 3. Checks if the provided password matches the stored password.
// 4. Generates access and refresh tokens for the authenticated user.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userLoginData (userModel.UserLogin): The data for user login.
//
// Returns:
// - commonModel.Result[userModel.UserToken]: The result containing either the user token data or an error.
func (userUseCaseV1 UserUseCaseV1) Login(ctx context.Context, userLoginData userModel.UserLogin) commonModel.Result[userModel.UserToken] {
	// Validate the user login data.
	userLogin := validateUserLogin(userLoginData)
	if validator.IsError(userLogin.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userLogin.Error))
	}

	// Call the user repository to fetch the user by email.
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
// It performs the following steps:
// 1. Generates the UserTokenPayload using the user ID and role.
// 2. Calls the token generation utility to create new access and refresh tokens.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - user (userModel.User): The user information for generating the tokens.
//
// Returns:
// - commonModel.Result[userModel.UserToken]: The result containing either the user token data or an error.
func (userUseCaseV1 UserUseCaseV1) RefreshAccessToken(ctx context.Context, user userModel.User) commonModel.Result[userModel.UserToken] {
	// Generate the UserTokenPayload using the user ID and role.
	userTokenPayload := domainModel.NewUserTokenPayload(user.ID, user.Role)

	// Generate new access and refresh tokens.
	userToken := generateToken(userTokenPayload)

	// Return the result containing the updated token information.
	return userToken
}

// ForgottenPassword handles the process of initiating a password reset for a user.
// It performs the following steps:
// 1. Validates the provided user forgotten password data.
// 2. Calls the user repository to fetch the user by email.
// 3. Generates a password reset token and sets the reset properties.
// 4. Updates the user in the repository with the reset token and timestamp.
// 5. Prepares email data and sends the password reset email.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userForgottenPasswordData (userModel.UserForgottenPassword): The data for initiating the password reset.
//
// Returns:
// - error: The error, if any, encountered during the process.
func (userUseCaseV1 UserUseCaseV1) ForgottenPassword(ctx context.Context, userForgottenPasswordData userModel.UserForgottenPassword) error {
	// Validate the user forgotten password data.
	userForgottenPassword := validateUserForgottenPassword(userForgottenPasswordData)
	if validator.IsError(userForgottenPassword.Error) {
		return domainError.HandleError(userForgottenPassword.Error)
	}

	// Call the user repository to fetch the user by email.
	fetchedUser := userUseCaseV1.GetUserByEmail(ctx, userForgottenPassword.Data.Email)
	if validator.IsError(fetchedUser.Error) {
		return domainError.HandleError(fetchedUser.Error)
	}

	// Generate the password reset token and set reset properties.
	token := randstr.String(resetTokenLength)
	encodedToken := commonUtility.Encode(token)
	tokenExpirationTime := time.Now().Add(constants.PasswordResetTokenExpirationTime)
	userForgottenPassword.Data.ResetToken = token
	userForgottenPassword.Data.ResetExpiry = tokenExpirationTime

	// Update the user in the repository with the reset token and timestamp.
	updatedUserPasswordError := userUseCaseV1.userRepository.ForgottenPassword(ctx, userForgottenPassword.Data)
	if validator.IsError(updatedUserPasswordError) {
		return domainError.HandleError(updatedUserPasswordError)
	}

	// Prepare email data for user forgotten password and send the password reset email.
	emailData := prepareEmailDataForUpdatePasswordResetToken(fetchedUser.Data.Name, encodedToken)
	sendEmailError := userUseCaseV1.userRepository.SendEmail(fetchedUser.Data, emailData)
	if validator.IsError(sendEmailError) {
		return domainError.HandleError(sendEmailError)
	}

	// Return nil to indicate success.
	return nil
}

// ResetUserPassword handles the process of resetting a user's password based on a valid reset token.
// It performs the following steps:
// 1. Validates the provided user reset password data.
// 2. Calls the user repository to fetch the user by email.
// 3. Checks the validity of the provided reset token.
// 4. Updates the user's password and clears the reset token and timestamp in the repository.
//
// Parameters:
// - ctx (context.Context): The context for managing request deadlines and cancellation signals.
// - userResetPasswordData (userModel.UserResetPassword): The data for resetting the user's password.
//
// Returns:
// - error: The error, if any, encountered during the process.
func (userUseCaseV1 UserUseCaseV1) ResetUserPassword(ctx context.Context, userResetPasswordData userModel.UserResetPassword) error {
	// Decode the reset token.
	token := commonUtility.Decode(location, userResetPasswordData.ResetToken)
	if validator.IsError(token.Error) {
		return domainError.HandleError(token.Error)
	}

	userResetPasswordData.ResetToken = token.Data

	// Validate the user reset password data.
	userResetPassword := validateResetPassword(userResetPasswordData)
	if validator.IsError(userResetPassword.Error) {
		return domainError.HandleError(userResetPassword.Error)
	}

	// Fetch the user by the reset token.
	fetchedUser := userUseCaseV1.userRepository.GetUserByResetToken(ctx, token.Data)
	if validator.IsError(fetchedUser.Error) {
		return domainError.HandleError(fetchedUser.Error)
	}

	// Check if the reset token has expired.
	if validator.IsTimeNotValid(fetchedUser.Data.ResetExpiry) {
		timeExpiredError := domainError.NewTimeExpiredError(location, constants.TimeExpiredErrorNotification)
		logging.Logger(timeExpiredError)
		return domainError.HandleError(timeExpiredError)
	}

	// Update the user's password.
	resetUserPasswordError := userUseCaseV1.userRepository.ResetUserPassword(ctx, userResetPassword.Data)
	if validator.IsError(resetUserPasswordError) {
		return domainError.HandleError(resetUserPasswordError)
	}

	// Return nil to indicate success.
	return nil
}

// prepareEmailData is a helper function to create an EmailData model for sending an email.
// It takes the user name, token value, email subject, URL, template name, and template path as input.
// It constructs an EmailData model and returns it.
//
// Parameters:
// - userName (string): The name of the user to whom the email is sent.
// - tokenValue (string): The token to be included in the email.
// - subject (string): The subject of the email.
// - url (string): The URL to be included in the email.
// - templateName (string): The name of the email template.
// - templatePath (string): The path to the email template.
//
// Returns:
// - userModel.EmailData: The prepared email data.
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

// prepareEmailDataForRegister is a helper function to prepare an EmailData model specifically for user registration.
// It takes the user name and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
//
// Parameters:
// - userName (string): The name of the user to whom the email is sent.
// - tokenValue (string): The token to be included in the email.
//
// Returns:
// - userModel.EmailData: The prepared email data for user registration.
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

// prepareEmailDataForUserUpdate is a helper function to prepare an EmailData model specifically for updating user information.
// It takes the user name and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
//
// Parameters:
// - userName (string): The name of the user to whom the email is sent.
// - tokenValue (string): The token to be included in the email.
//
// Returns:
// - userModel.EmailData: The prepared email data for updating user information.
func prepareEmailDataForUpdatePasswordResetToken(userName, tokenValue string) userModel.EmailData {
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

// generateToken generates access and refresh tokens for a user.
// It takes the user token payload as input and uses the application configuration
// to create both access and refresh tokens using domainUtility.GenerateJWTToken.
// If there are errors during token generation, it returns a failure result
// with the corresponding error. Otherwise, it returns a success result
// containing the generated access and refresh tokens.
//
// Parameters:
// - userTokenPayload (domainModel.UserTokenPayload): The payload containing user information for token generation.
//
// Returns:
// - commonModel.Result[userModel.UserToken]: The result containing either the generated tokens or an error.
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
	if validator.IsError(refreshToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](refreshToken.Error)
	}

	// Update the userToken struct with the generated tokens.
	userToken.AccessToken = accessToken.Data
	userToken.RefreshToken = refreshToken.Data

	// Return a success result with the userToken struct.
	return commonModel.NewResultOnSuccess[userModel.UserToken](userToken)
}
