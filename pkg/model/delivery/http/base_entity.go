package http

import (
	"time"

	utility "github.com/yachnytskyi/golang-mongo-grpc/pkg/utility/common"
)

type BaseEntity struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

func NewBaseEntity(id string, createdAt, updatedAt time.Time) BaseEntity {
	createdAtFormatted := utility.FormatDate(createdAt)
	updatedAtFormatted := utility.FormatDate(updatedAt)
	return BaseEntity{
		ID:        id,
		CreatedAt: createdAtFormatted,
		UpdatedAt: updatedAtFormatted,
	}
}
