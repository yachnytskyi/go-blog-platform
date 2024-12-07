package model

import (
	"time"

	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/delivery/http"
)

type TokenView struct {
	Token string `json:"token"`
}

type UsersView struct {
	UsersView              []UserView                   `json:"users"`
	HTTPPaginationResponse model.HTTPPaginationResponse `json:"pagination_response"`
}

type UserView struct {
	model.BaseEntity
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

type UserCreateView struct {
	Username        string `json:"username"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserUpdateView struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type UserLoginView struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserTokenView struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type UserForgottenPasswordView struct {
	Email string `json:"email"`
}

type UserResetPasswordView struct {
	ResetToken      string `json:"reset_token"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

type UserWelcomeMessageView struct {
	Notification string `json:"notification"`
}

func NewWelcomeMessageView(notification string) UserWelcomeMessageView {
	return UserWelcomeMessageView{
		Notification: notification,
	}
}

func NewUsersView(users []UserView, paginationResponse model.HTTPPaginationResponse) UsersView {
	return UsersView{
		UsersView:              users,
		HTTPPaginationResponse: paginationResponse,
	}
}

func NewUserView(id string, username, email, role string, createdAt, updatedAt time.Time) UserView {
	return UserView{
		BaseEntity: model.NewBaseEntity(id, createdAt, updatedAt),
		Username:   username,
		Email:      email,
		Role:       role,
	}
}

func NewUserCreateView(username, email, password, passwordConfirm string) UserCreateView {
	return UserCreateView{
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

func NewUserUpdateView(id, username string) UserUpdateView {
	return UserUpdateView{
		ID:       id,
		Username: username,
	}
}

func NewUserLoginView(email, password string) UserLoginView {
	return UserLoginView{
		Email:    email,
		Password: password,
	}
}

func NewUserForgottenPasswordView(email string) UserForgottenPasswordView {
	return UserForgottenPasswordView{
		Email: email,
	}
}

func NewUserResetPasswordView(resetToken, password, passwordConfirm string) UserResetPasswordView {
	return UserResetPasswordView{
		ResetToken:      resetToken,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

func NewUserTokenView(accessToken, refreshToken string) UserTokenView {
	return UserTokenView{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}
