package model

import "time"

// [GET].
type PostsView struct {
	PostsView []*PostView `json:"posts"`
}

// [GET].
type PostView struct {
	PostID    string    `json:"post_id,omitempty" bson:"_id,omitempty" db:"post_id,omitempty"`
	Title     string    `json:"title,omitempty" bson:"title,omitempty" db:"title,omitempty"`
	Content   string    `json:"content,omitempty" bson:"content,omitempty" db:"content,omitempty"`
	Image     string    `json:"image,omitempty" bson:"image,omitempty" db:"image,omitempty"`
	User      string    `json:"user,omitempty" bson:"user,omitempty" db:"user,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty" bson:"created_at,omitempty" db:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty" bson:"updated_at,omitempty" db:"updated_at,omitempty"`
}
