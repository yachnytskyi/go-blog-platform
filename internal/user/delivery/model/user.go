package model

import "time"

// [GET].
type UsersView struct {
	UsersView []*UserView `json:"users"`
}

// [GET].
type UserView struct {
	Name      string    `json:"name" bson:"name" db:"name"`
	Email     string    `json:"email" bson:"email" db:"email"`
	Role      string    `json:"role" bson:"role" db:"role"`
	CreatedAt time.Time `json:"created_at" bson:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at" db:"updated_at"`
}
