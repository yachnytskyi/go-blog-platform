package model

import "time"

// [GET].
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

// [POST].
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

// [PUT].
type UserUpdate struct {
	Name      string
	UpdatedAt time.Time
}

// [POST].
type UserSignIn struct {
	Email    string `json:"email" bson:"email" db:"email" binding:"required,lte=40,email"`
	Password string `json:"password" bson:"password" db:"password" binding:"required,min=8"`
}

// [GET].
type UserForgottenPassword struct {
	Email string `json:"email" binding:"required"`
}

// [POST].
type UserResetPassword struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}
