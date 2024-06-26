package mongo

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// BaseEntity represents the base entity structure used in the MongoDB layer.
// It contains common fields like ID, CreatedAt, and UpdatedAt.
//
// Fields:
// - ID: The unique identifier of the entity, represented as a primitive.ObjectID for MongoDB.
// - CreatedAt: The timestamp when the entity was created.
// - UpdatedAt: The timestamp when the entity was last updated.
type BaseEntity struct {
	ID        primitive.ObjectID `bson:"_id"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

// NewBaseEntity creates a new instance of BaseEntity.
// Parameters:
// - id: The unique identifier of the entity, represented as a primitive.ObjectID for MongoDB.
// - createdAt: The timestamp when the entity was created.
// - updatedAt: The timestamp when the entity was last updated.
//
// Returns:
// - A new BaseEntity instance.
func NewBaseEntity(id primitive.ObjectID, createdAt, updatedAt time.Time) BaseEntity {
	return BaseEntity{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
