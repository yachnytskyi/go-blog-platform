package service

import (
	"context"
	"fmt"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
	"github.com/yachnytskyi/golang-mongo-grpc/pkg/utils"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type UserService struct {
	userRepository user.Repository
}

func NewUserService(userRepository user.Repository) user.Service {
	return &UserService{userRepository: userRepository}
}

func (userService *UserService) Register(ctx context.Context, user *models.UserCreate) (*models.User, error) {
	createdUser, err := userService.userRepository.Register(ctx, user)

	return createdUser, err
}

func (userService *UserService) Login(ctx context.Context, user *models.UserSignIn) (*models.User, error) {
	fetchedUser, err := userService.userRepository.GetUserByEmail(ctx, user.Email)

	// Will return wrong email or password.
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Verify password - we previously created this method.
	matchPasswords := utils.VerifyPassword(fetchedUser.Password, user.Password)

	if matchPasswords != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	return fetchedUser, err
}

func (userService *UserService) UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.User, error) {
	updatedUser, err := userService.userRepository.UpdateNewRegisteredUserById(ctx, userID, key, value)

	return updatedUser, err
}

func (userService *UserService) UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string,
	secondKey string, secondValue time.Time) error {
	updatedUser := userService.userRepository.UpdatePasswordResetTokenUserByEmail(ctx, email, firstKey, firstValue, secondKey, secondValue)

	return updatedUser
}

func (userService *UserService) ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error {
	updatedUser := userService.userRepository.ResetUserPassword(ctx, firstKey, firstValue, secondKey, passwordKey, password)

	return updatedUser

}

func (userService *UserService) UpdateUserById(ctx context.Context, userID string, user *models.UserUpdate) (*models.UserView, error) {
	updatedUser, err := userService.userRepository.UpdateUserById(ctx, userID, user)

	return updatedUser, err
}

func (userService *UserService) GetUserById(ctx context.Context, userID string) (*models.User, error) {
	fetchedUser, err := userService.userRepository.GetUserById(ctx, userID)

	return fetchedUser, err
}

func (userService *UserService) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	fetchedUser, err := userService.userRepository.GetUserByEmail(ctx, email)

	return fetchedUser, err
}

func (userService *UserService) DeleteUserById(ctx context.Context, userID string) error {
	deletedUser := userService.userRepository.DeleteUserById(ctx, userID)

	return deletedUser
}
