package model

import "time"

// [GET].
type UsersView struct {
	UsersView []*UserView `json:"users"`
}

// [GET].
type UserView struct {
	UserID    string    `json:"user_id,omitempty"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// [POST].
type UserCreateView struct {
	Name            string    `json:"name" binding:"required,min=5,lte=40"`
	Email           string    `json:"email" binding:"required,min=3,lte=40,email"`
	Password        string    `json:"password" binding:"required,min=8,lte=100"`
	PasswordConfirm string    `json:"password_confirm" binding:"required,min=8,lte=100"`
	Role            string    `json:"role"`
	Verified        bool      `json:"verifyed"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
