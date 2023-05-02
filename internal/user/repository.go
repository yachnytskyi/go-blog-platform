package user

import (
	"context"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Repository interface {
	GetUserById(ctx context.Context, userID string) (*models.UserFullResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserFullResponse, error)
	Register(ctx context.Context, user *models.UserCreate) (*models.UserFullResponse, error)
	UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.UserFullResponse, error)
	UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string, secondKey string, secondValue time.Time) error
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error
	UpdateUserById(ctx context.Context, userID string, user *models.UserUpdateDomain) (*models.UserResponse, error)
	DeleteUserById(ctx context.Context, userID string) error
}
