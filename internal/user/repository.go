package user

import (
	"context"
	"time"

	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	domainError "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/error/domain_error"
)

type Repository interface {
	GetAllUsers(ctx context.Context, page int, limit int) (*userModel.Users, error)
	GetUserById(ctx context.Context, userID string) (*userModel.User, error)
	GetUserByEmail(ctx context.Context, email string) (*userModel.User, error)
	Register(ctx context.Context, user *userModel.UserCreate) (*userModel.User, domainError.DomainError)
	UpdateNewRegisteredUserById(ctx context.Context, userID string, key string, value string) (*userModel.User, error)
	UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string, secondKey string, secondValue time.Time) error
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error
	UpdateUserById(ctx context.Context, userID string, user *userModel.UserUpdate) (*userModel.User, error)
	DeleteUserById(ctx context.Context, userID string) error
}
