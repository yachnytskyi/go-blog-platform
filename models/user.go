package models

import (
	"time"
)

// [POST].
type UserCreate struct {
	Name            string    `json:"name" bson:"name" db:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" db:"emal" binding:"required,lte=40,email"`
	Password        string    `json:"password" bson:"password" db:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"password_confirm" bson:"password_confirm,omitempty" db:"password_confirm,omitempty" binding:"required"`
	Role            string    `json:"role" bson:"role" db:"role"`
	Verified        bool      `json:"verified" bson:"verified" db:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

// [POST].
type UserCreateRepository struct {
	Name            string    `json:"name" bson:"name" db:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" db:"emal" binding:"required,lte=40,email"`
	Password        string    `json:"password" bson:"password" db:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"password_confirm" bson:"password_confirm,omitempty" db:"password_confirm,omitempty" binding:"required"`
	Role            string    `json:"role" bson:"role" db:"role"`
	Verified        bool      `json:"verified" bson:"verified" db:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

// [POST].
type UserSignIn struct {
	Email    string `json:"email" bson:"email" db:"email" binding:"required,lte=40,email"`
	Password string `json:"password" bson:"password" db:"password" binding:"required,min=8"`
}

// [GET].
type ForgottenPasswordInput struct {
	Email string `json:"email" binding:"required"`
}

// [POST].
type ResetUserPasswordInput struct {
	Password        string `json:"password" binding:"required"`
	PasswordConfirm string `json:"passwordConfirm" binding:"required"`
}

// [PUT].
type UserUpdate struct {
	Name      string    `json:"name" bson:"name" db:"name" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

// [PUT].
type UserUpdateRepository struct {
	Name      string    `json:"name" bson:"name" db:"name" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

// [GET].
type User struct {
	UserID          string    `json:"user_id" bson:"_id" db:"user_id"`
	Name            string    `json:"name" bson:"name" db:"name"`
	Email           string    `json:"email" bson:"email" db:"email"`
	Password        string    `json:"password" bson:"password" db:"password"`
	PasswordConfirm string    `json:"password_confirm" bson:"password_confirm,omitempty" db:"password_confirm,omitempty"`
	Role            string    `json:"role" bson:"role" db:"role"`
	Verified        bool      `json:"verified" bson:"verified" db:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

// [GET].
type UserView struct {
	UserID    string    `json:"user_id" bson:"_id" db:"user_id"`
	Name      string    `json:"name" bson:"name" db:"name"`
	Email     string    `json:"email" bson:"email" db:"email"`
	Role      string    `json:"role" bson:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

func UserToUserViewMapper(user *User) UserView {
	return UserView{
		UserID:    user.UserID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func UserCreateToUserCreateRepositoryMapper(user *UserCreate) UserCreateRepository {
	return UserCreateRepository{
		Name:            user.Name,
		Email:           user.Email,
		Password:        user.Password,
		PasswordConfirm: user.PasswordConfirm,
		Role:            user.Role,
		Verified:        user.Verified,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}
}

func UserUpdateToUserUpdateRepositoryMapper(user *UserUpdate) UserUpdateRepository {
	return UserUpdateRepository{
		Name:      user.Name,
		UpdatedAt: user.UpdatedAt,
	}
}
