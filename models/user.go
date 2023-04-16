package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// [POST].
type UserCreate struct {
	Name            string    `json:"name" bson:"name" binding:"required"`
	Email           string    `json:"email" bson:"email" binding:"required,lte=40,email"`
	Password        string    `json:"password" bson:"password" binding:"required,min=8"`
	PasswordConfirm string    `json:"password_confirm" bson:"password_confirm,omitempty" binding:"required"`
	Role            string    `json:"role" bson:"role"`
	Verified        bool      `json:"verified" bson:"verified"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}

// [POST].
type UserSignIn struct {
	Email    string `json:"email" bson:"email" binding:"required,lte=40,email"`
	Password string `json:"password" bson:"email" binding:"required,min=8"`
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
	Name      string    `json:"name" bson:"name" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

// [GET].
type UserDBFullResponse struct {
	UserID          primitive.ObjectID `json:"user_id" bson:"_id"`
	Name            string             `json:"name" bson:"name"`
	Email           string             `json:"email" bson:"email"`
	Password        string             `json:"password" bson:"password"`
	PasswordConfirm string             `json:"password_confirm" bson:"password_confirm,omitempty"`
	Role            string             `json:"role" bson:"role"`
	Verified        bool               `json:"verified" bson:"verified"`
	CreatedAt       time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at" bson:"updated_at"`
}

// [GET].
type UserResponse struct {
	UserID    primitive.ObjectID `json:"user_id" bson:"_id"`
	Name      string             `json:"name" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Role      string             `json:"role" bson:"role"`
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time          `json:"updated_at" bson:"updated_at"`
}

// User mapping.
func FilteredResponse(user *UserDBFullResponse) UserResponse {
	return UserResponse{
		UserID:    user.UserID,
		Email:     user.Email,
		Name:      user.Name,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
