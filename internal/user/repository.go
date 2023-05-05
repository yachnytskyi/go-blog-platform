package user

import (
	"context"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Repository interface {
	GetUserById(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	Register(ctx context.Context, user *models.UserCreateDomain) (*models.User, error)
	UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.User, error)
	UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string, secondKey string, secondValue time.Time) error
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error
	UpdateUserById(ctx context.Context, userID string, user *models.UserUpdateDomain) (*models.UserView, error)
	DeleteUserById(ctx context.Context, userID string) error
}
