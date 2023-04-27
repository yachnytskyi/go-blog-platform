package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// [POST].
type PostCreate struct {
	Title     string    `json:"title" bson:"title" binding:"required"`
	Content   string    `json:"content" bson:"content" binding:"required"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty"`
	UserID    string    `json:"user_id,omitempty" bson:"user_id,omitempty"`
	User      string    `json:"user" bson:"user" binding:"required"`
	CreateAt  time.Time `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// [GET].
type PostDB struct {
	PostID    primitive.ObjectID `json:"post_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty"`
	UserID    string             `json:"-" bson:"user_id"`
	User      string             `json:"user,omitempty" bson:"user,omitempty"`
	CreateAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}

// [PUT].
type PostUpdate struct {
	PostID    primitive.ObjectID `json:"post_id,omitempty" bson:"_id,omitempty"`
	Title     string             `json:"title,omitempty" bson:"title,omitempty"`
	Content   string             `json:"content,omitempty" bson:"content,omitempty"`
	Image     string             `json:"image,omitempty" bson:"image,omitempty"`
	UserID    string             `json:"-,omitempty" bson:"user_id,omitempty"`
	User      string             `json:"user,omitempty" bson:"user,omitempty"`
	CreateAt  time.Time          `json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty" bson:"updated_at,omitempty"`
}
