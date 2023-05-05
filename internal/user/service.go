package user

import (
	"context"
	"time"

	"github.com/yachnytskyi/golang-mongo-grpc/models"
)

type Service interface {
	GetUserById(ctx context.Context, userID string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error)
	Register(ctx context.Context, user *models.UserCreateDomain) (*models.User, error)
	Login(ctx context.Context, user *models.UserSignIn) (*models.User, error)
	UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*models.User, error)
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey, passwordKey, password string) error
	UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string, secondKey string, secondValue time.Time) error
	UpdateUserById(ctx context.Context, userID string, user *models.UserUpdateDomain) (*models.UserView, error)
	DeleteUserById(ctx context.Context, userID string) error
}
