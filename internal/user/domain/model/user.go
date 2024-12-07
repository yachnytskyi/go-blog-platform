package model

import (
	"time"

	common "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/common"
	model "github.com/yachnytskyi/golang-mongo-grpc/pkg/model/domain"
)

type Users struct {
	Users              []User
	PaginationResponse common.PaginationResponse
}

type User struct {
	model.BaseEntity
	Username string
	Email    string
	Password string
	Role     string
	Verified bool
}

type UserCreate struct {
	Username         string
	Email            string
	Password         string
	PasswordConfirm  string
	Role             string
	Verified         bool
	VerificationCode string
}

type UserUpdate struct {
	ID        string
	Username  string
}

type UserLogin struct {
	Email    string
	Password string
}

type UserToken struct {
	AccessToken  string
	RefreshToken string
}

type UserForgottenPassword struct {
	Email       string
	ResetToken  string
	ResetExpiry time.Time
}

type UserResetPassword struct {
	ResetToken      string
	Password        string
	PasswordConfirm string
}

type UserResetExpiry struct {
	ResetExpiry time.Time
}

func NewUsers(users []User, paginationResponse common.PaginationResponse) Users {
	return Users{
		Users:              users,
		PaginationResponse: paginationResponse,
	}
}

func NewUser(id string, username, email, password, role string, verified bool, createdAt, updatedAt time.Time) User {
	return User{
		BaseEntity: model.NewBaseEntity(id, createdAt, updatedAt),
		Username:   username,
		Email:      email,
		Password:   password,
		Role:       role,
		Verified:   verified,
	}
}

func NewUserCreate(username, email, password, passwordConfirm string) UserCreate {
	return UserCreate{
		Username:        username,
		Email:           email,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

func NewUserUpdate(id, username string) UserUpdate {
	return UserUpdate{
		ID:       id,
		Username: username,
	}
}

func NewUserLogin(email, password string) UserLogin {
	return UserLogin{
		Email:    email,
		Password: password,
	}
}

func NewUserToken(accessToken, refreshToken string) UserToken {
	return UserToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
}

func NewUserForgottenPassword(email string) UserForgottenPassword {
	return UserForgottenPassword{
		Email: email,
	}
}

func NewUserResetPassword(resetToken, password, passwordConfirm string) UserResetPassword {
	return UserResetPassword{
		ResetToken:      resetToken,
		Password:        password,
		PasswordConfirm: passwordConfirm,
	}
}

func NewUserResetExpiry(resetExpiry time.Time) UserResetExpiry {
	return UserResetExpiry{
		ResetExpiry: resetExpiry,
	}
}
