package models

import (
	"time"
)

// [POST].
type PostCreate struct {
	Title     string    `json:"title" bson:"title" db:"title" binding:"required"`
	Content   string    `json:"content" bson:"content" db:"content" binding:"required"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty" db:"image,omitempty"`
	UserID    string    `json:"user_id" bson:"user_id" db:"user_id" binding:"required"`
	User      string    `json:"user" bson:"user" db:"user" binding:"required"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at,omiptempty" bson:"updated_at,omitempty"`
}

// [GET].
type PostDB struct {
	PostID    string    `json:"post_id,omitempty" bson:"_id,omitempty" db:"post_id,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty"`
	Content   string    `json:"content,omitempty" bson:"content,omitempty" db:"content,omitempty"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty" db:"image,omitempty"`
	UserID    string    `json:"-" bson:"user_id" db:"user_id" binding:"required"`
	User      string    `json:"user,omitempty" bson:"user,omitempty" db:"user,omitpempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty" db:"updated_at,omitempty"`
}

// [PUT].
type PostUpdate struct {
	PostID    string    `json:"post_id,omitempty" bson:"_id,omitempty" db:"post_id,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty" db:"title,omitempty"`
	Content   string    `json:"content,omitempty" bson:"content,omitempty"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty" db:",image"`
	UserID    string    `json:"user_id" bson:"user_id" db:"user_id"`
	User      string    `json:"user,omitempty" bson:"user,omitempty" db:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty" db:"updated_at,omitempty"`
}
