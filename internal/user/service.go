package user

import (
	"context"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Service interface {
	Register(ctx context.Context, user *models.UserCreate) (*models.UserDB, error)
	Login(ctx context.Context, user *models.UserSignIn) (*models.UserDB, error)
	UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.UserDB, error)
	UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string, secondKey string, secondValue time.Time) error
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey, passwordKey, password string) error
	UpdateUserById(ctx context.Context, userID string, user *models.UserUpdate) (*models.UserDB, error)
	GetUserById(ctx context.Context, userID string) (*models.UserDB, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserDB, error)
	DeleteUserById(ctx context.Context, userID string) error
}
