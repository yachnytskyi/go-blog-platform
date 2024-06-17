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
	tokenValue := randstr.String(verificationCodeLength)
	tokenValue = commonUtility.Encode(tokenValue)
	userCreate.Data.Role = userRole
	userCreate.Data.Verified = true
	userCreate.Data.VerificationCode = tokenValue
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
	emailData := prepareEmailDataForRegistration(createdUser.Data.Name, tokenValue)
	sendEmailVerificationMessageError := userUseCaseV1.userRepository.SendEmailVerificationMessage(ctx, createdUser.Data, emailData)
	if validator.IsError(sendEmailVerificationMessageError) {
		return commonModel.NewResultOnFailure[userModel.User](domainError.HandleError(sendEmailVerificationMessageError))
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
	userTokenPayload := domainModel.NewUserTokenPayload(fetchedUser.Data.UserID, fetchedUser.Data.Role)

	// Generate access and refresh tokens.
	userToken := generateToken(userTokenPayload)
	if validator.IsError(userLogin.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](domainError.HandleError(userToken.Error))
	}

	// Return the result containing userToken information.
	return userToken
}

func (userUseCaseV1 UserUseCaseV1) RefreshAccessToken(ctx context.Context, user userModel.User) commonModel.Result[userModel.UserToken] {
	// Generate the UserTokenPayload.
	userTokenPayload := domainModel.NewUserTokenPayload(user.UserID, user.Role)

	// Generate access and refresh tokens.
	userToken := generateToken(userTokenPayload)

	// Return the result containing userToken information.
	return userToken
}

func (userUseCaseV1 UserUseCaseV1) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {
	validateEmailError := isEmailValid(email)
	if validator.IsError(validateEmailError) {
		return domainError.HandleError(validateEmailError)
	}
	updatedUserError := userUseCaseV1.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, email, firstKey, firstValue, secondKey, secondValue)
	if validator.IsError(updatedUserError) {
		updatedUserError = domainError.HandleError(updatedUserError)
		return updatedUserError
	}

	// Generate verification code.
	tokenValue := randstr.String(resetTokenLength)
	encodedTokenValue := commonUtility.Encode(tokenValue)
	tokenExpirationTime := time.Now().Add(time.Minute * 15)

	// Update the user.
	fetchedUser := userUseCaseV1.GetUserByEmail(ctx, email)
	if validator.IsError(fetchedUser.Error) {
		fetchedUserError := domainError.HandleError(fetchedUser.Error)
		return fetchedUserError
	}
	updatedUserPasswordError := userUseCaseV1.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, fetchedUser.Data.Email, "passwordResetToken", encodedTokenValue, "passwordResetAt", tokenExpirationTime)
	if validator.IsError(updatedUserPasswordError) {
		updatedUserPasswordError = domainError.HandleError(updatedUserPasswordError)
		return updatedUserPasswordError
	}

	emailData := prepareEmailDataForUpdatePasswordResetToken(fetchedUser.Data.Name, tokenValue)
	sendEmailForgottenPasswordMessageError := userUseCaseV1.userRepository.SendEmailForgottenPasswordMessage(ctx, fetchedUser.Data, emailData)
	if validator.IsError(sendEmailForgottenPasswordMessageError) {
		sendEmailForgottenPasswordMessageError = domainError.HandleError(sendEmailForgottenPasswordMessageError)
		return sendEmailForgottenPasswordMessageError
	}

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
	emailConfig := config.AppConfig.Email
	userFirstName := domainUtility.UserFirstName(userName)
	emailData := userModel.NewEmailData(emailConfig.ClientOriginUrl+url+tokenValue, templateName, templatePath, userFirstName, subject)
	return emailData
}

// prepareEmailDataForRegister is a helper function to prepare an EmailData model specifically for user registration.
// It takes the context, user name, and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
func prepareEmailDataForRegistration(userName, tokenValue string) userModel.EmailData {
	emailConfig := config.AppConfig.Email
	return prepareEmailData(userName, tokenValue, constants.EmailConfirmationSubject, constants.EmailConfirmationUrl,
		emailConfig.UserConfirmationTemplateName, emailConfig.UserConfirmationTemplatePath)
}

// prepareEmailDataForUserUpdate is a helper function to prepare an EmailData model specifically for updating user information.
// It takes the context, user name, and token value as input and uses the constants for email subject, URL, template name, and template path.
// It internally calls prepareEmailData with the appropriate parameters and returns the result.
func prepareEmailDataForUpdatePasswordResetToken(userName, tokenValue string) userModel.EmailData {
	emailConfig := config.AppConfig.Email
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
	accessTokenConfig := config.AppConfig.AccessToken
	refreshTokenConfig := config.AppConfig.RefreshToken

	// Generate the access token.
	accessToken := domainUtility.GenerateJWTToken(location, accessTokenConfig.PrivateKey, accessTokenConfig.ExpiredIn, userTokenPayload)
	if validator.IsError(accessToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](accessToken.Error)
	}

	// Generate the refresh token.
	refreshToken := domainUtility.GenerateJWTToken(location, refreshTokenConfig.PrivateKey, refreshTokenConfig.ExpiredIn, userTokenPayload)
	if validator.IsError(accessToken.Error) {
		return commonModel.NewResultOnFailure[userModel.UserToken](refreshToken.Error)
	}

	// Update the userToken struct with the generated tokens.
	userToken.AccessToken = accessToken.Data
	userToken.RefreshToken = refreshToken.Data

	// Return a success result with the userToken struct.
	return commonModel.NewResultOnSuccess[userModel.UserToken](userToken)
}
