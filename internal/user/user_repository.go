package user

import (
	"context"

	userModel "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	commonModel "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context, paginationQuery commonModel.PaginationQuery) commonModel.Result[userModel.Users]
	GetUserById(ctx context.Context, userID string) commonModel.Result[userModel.User]
	GetUserByEmail(ctx context.Context, email string) commonModel.Result[userModel.User]
	CheckEmailDuplicate(ctx context.Context, email string) error
	SendEmail(user userModel.User, data userModel.EmailData) error
	Register(ctx context.Context, user userModel.UserCreate) commonModel.Result[userModel.User]
	UpdateCurrentUser(ctx context.Context, user userModel.UserUpdate) commonModel.Result[userModel.User]
	DeleteUserById(ctx context.Context, userID string) error
	ForgottenPassword(ctx context.Context, userForgottenPassword userModel.UserForgottenPassword) error
	ResetUserPassword(ctx context.Context, firstKey string, firstValue string, secondKey string, passwordKey, password string) error
}
