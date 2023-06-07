package model

import "time"

type Users struct {
	Users []*User
	Limit int
}

type User struct {
	UserID    string
	Name      string
	Email     string
	Password  string
	Role      string
	Verified  bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserCreate struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            string
	Verified        bool
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type UserUpdate struct {
	Name      string
	UpdatedAt time.Time
}

type UserLogin struct {
	Email    string
	Password string
}

type UserForgottenPassword struct {
	Email string `json:"email" binding:"required"`
}

type UserResetPassword struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}
