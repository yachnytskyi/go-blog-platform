package user

import (
	"context"

	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
)

type UserRepository interface {
	GetAllUsers(ctx context.Context, paginationQuery *common.PaginationQuery) common.Result[*user.Users]
	GetUserById(ctx context.Context, userID string) common.Result[user.User]
	GetUserByEmail(ctx context.Context, email string) common.Result[user.User]
	CheckEmailDuplicate(ctx context.Context, email string) error
	SendEmail(user user.User, data user.EmailData) error
	Register(ctx context.Context, user user.UserCreate) common.Result[user.User]
	UpdateCurrentUser(ctx context.Context, user user.UserUpdate) common.Result[user.User]
	DeleteUserById(ctx context.Context, userID string) error
	ForgottenPassword(ctx context.Context, userForgottenPassword user.UserForgottenPassword) error
	ResetUserPassword(ctx context.Context, userResetPassword user.UserResetPassword) error
	GetResetExpiry(ctx context.Context, token string) common.Result[user.UserResetExpiry]
}
