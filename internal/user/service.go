package user

import (
	"context"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Service interface {
	Register(ctx context.Context, user *models.UserCreate) (*models.UserDBFullResponse, error)
	Login(ctx context.Context, user *models.UserSignIn) (*models.UserDBFullResponse, error)
	UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.UserDBFullResponse, error)
	UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string, secondKey string, secondValue time.Time) error
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey, passwordKey, password string) error
	UpdateUserById(ctx context.Context, userID string, user *models.UserUpdate) (*models.UserDBFullResponse, error)
	GetUserById(ctx context.Context, userID string) (*models.UserDBFullResponse, error)
	GetUserByEmail(ctx context.Context, email string) (*models.UserDBFullResponse, error)
	DeleteUserById(ctx context.Context, userID string) error
}
