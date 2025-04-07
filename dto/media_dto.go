package dto

import (
	"time"

	"github.com/google/uuid"
)

// Book DTOs
type MediaRes struct {
	ID        uuid.UUID `json:"id"`
	Path      string    `json:"path"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
