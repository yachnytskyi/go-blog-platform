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
