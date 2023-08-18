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
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utility"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
)

type UserUseCase struct {
	userRepository user.Repository
}

func NewUserUseCase(userRepository user.Repository) user.UseCase {
	return &UserUseCase{userRepository: userRepository}
}

func (userUseCase *UserUseCase) GetAllUsers(ctx context.Context, page int, limit int) (*userModel.Users, error) {
	fetchedUsers, err := userUseCase.userRepository.GetAllUsers(ctx, page, limit)

	return fetchedUsers, err
}

func (userUseCase *UserUseCase) GetUserById(ctx context.Context, userID string) (*userModel.User, error) {
	fetchedUser, err := userUseCase.userRepository.GetUserById(ctx, userID)

	return fetchedUser, err
}

func (userUseCase *UserUseCase) GetUserByEmail(ctx context.Context, email string) (*userModel.User, error) {
	fetchedUser, err := userUseCase.userRepository.GetUserByEmail(ctx, email)

	return fetchedUser, err
}

func (userUseCase *UserUseCase) Register(ctx context.Context, user *userModel.UserCreate) (*userModel.User, []error) {
	if userUseCase.userRepository.CheckEmailDublicate(ctx, user.Email) {
		userCreateValidationErrors := []error{}

		userCreateValidationError := &domainError.ValidationError{
			Field:        "email",
			FieldType:    "required",
			Notification: EmailAlreadyExists,
		}

		userCreateValidationErrors = append(userCreateValidationErrors, userCreateValidationError)
		userCreateValidationErrors = domainError.ErrorsHandler(userCreateValidationErrors)
		return nil, userCreateValidationErrors
	}

	if userCreateValidationErrors := UserCreateValidator(user); len(userCreateValidationErrors) != 0 {
		userCreateValidationErrors = domainError.ErrorsHandler(userCreateValidationErrors)
		return nil, userCreateValidationErrors
	}

	createdUser, userCreateError := userUseCase.userRepository.Register(ctx, user)

	if userCreateError != nil {
		userCreateErrors := []error{userCreateError}
		userCreateErrors = domainError.ErrorsHandler(userCreateErrors)
		return nil, userCreateErrors
	}

	// Generate verification code.
	code := randstr.String(20)
	verificationCode := utility.Encode(code)

	_, userUpdateError := userUseCase.userRepository.UpdateNewRegisteredUserById(ctx, createdUser.UserID, "verificationCode", verificationCode)

	if userUpdateError != nil {
		userUpdateErrors := []error{userUpdateError}
		userUpdateErrors = domainError.ErrorsHandler(userUpdateErrors)
		return nil, userUpdateErrors
	}

	// Send an email.
	configLoad, err := config.LoadConfig(".")

	if err != nil {
		var sendEmailInternalError *domainError.InternalError = new(domainError.InternalError)
		sendEmailInternalError.Location = "User.Data.Repository.External.Mail.SendEmail.LoadConfig"
		sendEmailInternalError.Reason = err.Error()
		fmt.Println(sendEmailInternalError)
		sendEmailErrors := []error{sendEmailInternalError}
		sendEmailErrors = domainError.ErrorsHandler(sendEmailErrors)
		return nil, sendEmailErrors
	}

	firstName := domainUtility.UserFirstName(createdUser.Name)
	emailData := userModel.EmailData{
		URL:       configLoad.Origin + "/verifyemail/" + code,
		FirstName: firstName,
		Subject:   "Your account verification code",
	}

	if userUseCase.userRepository.SendEmailVerificationMessage(createdUser, &emailData, config.TemplateName) != nil {
		userCreateErrorMessages := []error{}
		userCreateErrorMessage := &domainError.ValidationError{
			Notification: domainError.InternalErrorNotification,
		}

		userCreateErrorMessages = append(userCreateErrorMessages, userCreateErrorMessage)
		userCreateErrorMessages = domainError.ErrorsHandler(userCreateErrorMessages)
		return nil, userCreateErrorMessages
	}

	return createdUser, nil
}

func (userUseCase *UserUseCase) UpdateUserById(ctx context.Context, userID string, user *userModel.UserUpdate) (*userModel.User, []error) {
	if userUpdateValidationErrors := UserUpdateValidator(user); len(userUpdateValidationErrors) != 0 {
		return nil, userUpdateValidationErrors
	}

	updatedUser, userUpdateError := userUseCase.userRepository.UpdateUserById(ctx, userID, user)

	if userUpdateError != (domainError.InternalError{}) {
		userUpdateErrors := []error{userUpdateError}
		return nil, userUpdateErrors
	}

	return updatedUser, []error{}
}

func (userUseCase *UserUseCase) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userUseCase.userRepository.DeleteUserById(ctx, userID)

	return deletedUser
}

func (userUseCase *UserUseCase) Login(ctx context.Context, user *userModel.UserLogin) (string, error) {
	if err := UserLoginValidator(user); err != nil {
		return "", err
	}

	fetchedUser, err := userUseCase.userRepository.GetUserByEmail(ctx, user.Email)

	// Will return wrong email or password.
	if err != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	// Verify password - we previously created this method.
	matchPasswords := domainUtility.VerifyPassword(fetchedUser.Password, user.Password)

	if matchPasswords != nil {
		return "", fmt.Errorf("invalid email or password")
	}

	return fetchedUser.UserID, err
}

func (userUseCase *UserUseCase) UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*userModel.User, error) {
	updatedUser, err := userUseCase.userRepository.UpdateNewRegisteredUserById(ctx, userID, key, value)

	return updatedUser, err
}

func (userUseCase *UserUseCase) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {
	updatedUser := userUseCase.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, email, firstKey, firstValue, secondKey, secondValue)

	return updatedUser
}

func (userUseCase *UserUseCase) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	updatedUser := userUseCase.userRepository.ResetUserPassword(ctx, firstKey, firstValue, secondKey, passwordKey, password)

	return updatedUser

}
