package service

import (
	"context"

	"github.com/yachnytskyi/golang-mongo-grpc/internal/user"
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

func (userService *UserService) UserGetById(ctx context.Context, userID string) (*models.UserFullResponse, error) {
	user, err := userService.userRepository.UserGetById(ctx, userID)

	return user, err
}

func (userService *UserService) UserGetByEmail(ctx context.Context, email string) (*models.UserFullResponse, error) {
	user, err := userService.userRepository.UserGetByEmail(ctx, email)

	return user, err
}
