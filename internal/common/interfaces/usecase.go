package interfaces

import (
	"context"

	post "github.com/yachnytskyi/golang-mongo-grpc/internal/post/domain/model"
	user "github.com/yachnytskyi/golang-mongo-grpc/internal/user/domain/model"
	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
)

type UserUseCase interface {
	GetAllUsers(ctx context.Context, paginationQuery common.PaginationQuery) common.Result[user.Users]
	GetUserById(ctx context.Context, userID string) common.Result[user.User]
	GetUserByEmail(ctx context.Context, email string) common.Result[user.User]
	Register(ctx context.Context, user user.UserCreate) common.Result[user.User]
	UpdateCurrentUser(ctx context.Context, user user.UserUpdate) common.Result[user.User]
	DeleteUserById(ctx context.Context, userID string) error
	Login(ctx context.Context, user user.UserLogin) common.Result[user.UserToken]
	RefreshAccessToken(ctx context.Context, user user.User) common.Result[user.UserToken]
	ForgottenPassword(ctx context.Context, userForgottenPasswordView user.UserForgottenPassword) error
	ResetUserPassword(ctx context.Context, userResetPassword user.UserResetPassword) error
}

type PostUseCase interface {
	GetAllPosts(ctx context.Context, page int, limit int) (*post.Posts, error)
	GetPostById(ctx context.Context, postID string) (*post.Post, error)
	CreatePost(ctx context.Context, user *post.PostCreate) (*post.Post, error)
	UpdatePostById(ctx context.Context, postID string, post *post.PostUpdate, currentUserID string) (*post.Post, error)
	DeletePostByID(ctx context.Context, postID string, currentUserID string) error
}
