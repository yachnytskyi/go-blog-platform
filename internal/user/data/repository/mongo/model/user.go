package model

import "time"

// [GET].
type UserRepository struct {
	UserID    string    `bson:"_id"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	Role      string    `bson:"role"`
	Verified  bool      `bson:"verified"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// [POST].
type UserCreateRepository struct {
	Name      string    `json:"name" bson:"name" db:"name" binding:"required"`
	Email     string    `json:"email" bson:"email" db:"emal" binding:"required,lte=40,email"`
	Password  string    `json:"password" bson:"password" db:"password" binding:"required,min=8"`
	Role      string    `json:"role" bson:"role" db:"role"`
	Verified  bool      `json:"verified" bson:"verified" db:"verified"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}

// [PUT].
type UserUpdateRepository struct {
	Name      string    `json:"name" bson:"name" db:"name" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}
