package service

import (
	"context"
	"fmt"

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

func (userService *UserService) Register(ctx context.Context, user *models.UserCreate) (*models.UserFullResponse, error) {
	createdUser, err := userService.userRepository.Register(ctx, user)

	return createdUser, err
}

func (userService *UserService) Login(ctx context.Context, user *models.UserSignIn) (*models.UserFullResponse, error) {
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

func (userService *UserService) GetUserById(ctx context.Context, userID string) (*models.UserFullResponse, error) {
	user, err := userService.userRepository.GetUserById(ctx, userID)

	return user, err
}

func (userService *UserService) GetUserByEmail(ctx context.Context, email string) (*models.UserFullResponse, error) {
	user, err := userService.userRepository.GetUserById(ctx, email)

	return user, err
}

func (userService *UserService) UpdateUserById(ctx context.Context, userID string, key string, value string) (*models.UserFullResponse, error) {
	user, err := userService.userRepository.UpdateUserById(ctx, userID, key, value)

	return user, err
}
