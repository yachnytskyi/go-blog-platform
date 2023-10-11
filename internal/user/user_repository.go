package user

import (
	"context"
	"time"

	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users]
	GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User]
	GetUserByEmail(ctx context.Context, email string) (userModel.User, error)
	CheckEmailDublicate(ctx context.Context, email string) error
	SendEmailVerificationMessage(ctx context.Context, user userModel.User, data userModel.EmailData) error
	SendEmailForgottenPasswordMessage(ctx context.Context, user userModel.User, data userModel.EmailData) error
	Register(ctx context.Context, user userModel.UserCreate) commonModel.Result[userModel.User]
	UpdateUserById(ctx context.Context, userID string, user userModel.UserUpdate) (userModel.User, error)
	DeleteUser(ctx context.Context, userID string) error
	UpdatePasswordResetTokenUserByEmail(ctx context.Context, email string, firstKey string, firstValue string, secondKey string, secondValue time.Time) error
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error
}
