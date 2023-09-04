package model

import (
	"time"
)

// [GET].
type UsersView struct {
	UsersView []*UserView `json:"users"`
}

// [GET].
type UserView struct {
	UserID    string    `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// [POST].
type UserCreateView struct {
	Name            string `json:"name"`
	Email           string `json:"email"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// [PUT].
type UserUpdateView struct {
	Name      string    `json:"name"`
	UpdatedAt time.Time `json:"updated_at"`
}

// [POST].
type UserLoginView struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// [POST].
type UserForgottenPasswordView struct {
	Email string `json:"email"`
}

// [POST].
type UserResetPasswordView struct {
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
}

// [GET].
type UserWelcomeMessageView struct {
	Message string `json:"message"`
}
