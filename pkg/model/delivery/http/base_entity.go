package http

import (
	"time"

	commonUtility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
)

// BaseEntity represents the base entity structure used in the HTTP layer.
// It encapsulates common fields such as ID, CreatedAt, and UpdatedAt which
// are shared across various entities in the HTTP context.
//
// Fields:
// - ID: The unique identifier of the entity, represented as a string.
// - CreatedAt: The timestamp indicating when the entity was created, formatted as a string.
// - UpdatedAt: The timestamp indicating the last update time of the entity, formatted as a string.
type BaseEntity struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

// NewBaseEntity creates a new instance of BaseEntity.
// It initializes the BaseEntity with the provided ID, createdAt, and updatedAt
// timestamps, formatting the timestamps using FormatDate function.
//
// Parameters:
// - id: The unique identifier of the entity, passed as a string.
// - createdAt: The time.Time value indicating when the entity was created.
// - updatedAt: The time.Time value indicating the last time the entity was updated.
//
// Returns:
// - A new instance of BaseEntity with the ID and formatted timestamps for CreatedAt and UpdatedAt.
func NewBaseEntity(id string, createdAt, updatedAt time.Time) BaseEntity {
	// Format the createdAt timestamp to a string using the common utility function
	createdAtFormatted := commonUtility.FormatDate(createdAt)

	// Format the updatedAt timestamp to a string using the common utility function
	updatedAtFormatted := commonUtility.FormatDate(updatedAt)

	// Return a new instance of BaseEntity with the provided ID and formatted timestamps
	return BaseEntity{
		ID:        id,
		CreatedAt: createdAtFormatted,
		UpdatedAt: updatedAtFormatted,
	}
}
