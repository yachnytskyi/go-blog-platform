package http

import "time"

// BaseEntity represents the base entity structure used in the HTTP layer.
// It contains common fields like ID, CreatedAt, and UpdatedAt.
//
// Fields:
// - ID: The unique identifier of the entity.
// - CreatedAt: The timestamp when the entity was created.
// - UpdatedAt: The timestamp when the entity was last updated.
type BaseEntity struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewBaseEntity creates a new instance of BaseEntity.
// Parameters:
// - id: The unique identifier of the entity.
// - createdAt: The timestamp when the entity was created.
// - updatedAt: The timestamp when the entity was last updated.
//
// Returns:
// - A new BaseEntity instance.
func NewBaseEntity(id string, createdAt, updatedAt time.Time) BaseEntity {
	return BaseEntity{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
