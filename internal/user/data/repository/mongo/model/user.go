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
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	Role      string    `bson:"role"`
	Verified  bool      `bson:"verified"`
	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

// [PUT].
type UserUpdateRepository struct {
	Name      string    `bson:"name"`
	UpdatedAt time.Time `bson:"updated_at"`
}
