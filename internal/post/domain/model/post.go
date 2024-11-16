package model

import "time"

type Posts struct {
	Posts []*Post
}

type Post struct {
	PostID    string
	UserID    string
	Title     string
	Content   string
	Image     string
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostCreate struct {
	UserID    string
	Title     string
	Content   string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type PostUpdate struct {
	PostID    string
	UserID    string
	Title     string
	Content   string
	Image     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
