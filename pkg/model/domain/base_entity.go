package domain

import "time"

type BaseEntity struct {
	ID        string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewBaseEntity(id string, createdAt, updatedAt time.Time) BaseEntity {
	return BaseEntity{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
