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
	Email string
}

type UserResetPassword struct {
	Password        string
	PasswordConfirm string
}

type EmailData struct {
	URL          string
	TemplateName string
	TemplatePath string
	FirstName    string
	Subject      string
}
