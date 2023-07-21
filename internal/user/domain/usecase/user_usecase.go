package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	domainUtility "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/utility"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"

	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
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
			Notification: userModel.EmailAlreadyExists,
		}

		userCreateValidationErrors = append(userCreateValidationErrors, userCreateValidationError)

		return nil, userCreateValidationErrors
	}

	if userCreateValidationErrors := user.UserCreateValidator(); len(userCreateValidationErrors) != 0 {
		return nil, userCreateValidationErrors
	}

	createdUser, userCreateError := userUseCase.userRepository.Register(ctx, user)

	if userCreateError != (domainError.InternalError{}) {
		userCreateErrors := []error{userCreateError}

		return nil, userCreateErrors
	}

	return createdUser, []error{}
}

func (userUseCase *UserUseCase) UpdateUserById(ctx context.Context, userID string, user *userModel.UserUpdate) (*userModel.User, []error) {
	if userUpdateValidationErrors := user.UserUpdateValidator(); len(userUpdateValidationErrors) != 0 {
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
	if err := user.UserLoginValidator(); err != nil {
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
