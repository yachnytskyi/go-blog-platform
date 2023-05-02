package user

import (
	"context"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Service interface {
	GetUserById(ctx context.Context, userID string) (*models.UserFullResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserFullResponse, error)
	Register(ctx context.Context, user *models.UserCreateDomain) (*models.UserFullResponse, error)
	Login(ctx context.Context, user *models.UserSignIn) (*models.UserFullResponse, error)
	UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.UserFullResponse, error)
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey, passwordKey, password string) error
	UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string, secondKey string, secondValue time.Time) error
	UpdateUserById(ctx context.Context, userID string, user *models.UserUpdateDomain) (*models.UserResponse, error)
	DeleteUserById(ctx context.Context, userID string) error
}
