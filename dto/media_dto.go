package dto

import (
	"go-gin-simple-api/model"
	"time"

	"github.com/google/uuid"
)

// Book DTOs
type MediaRes struct {
	ID        uuid.UUID    `json:"id"`
	Path      string       `json:"path"`
	Books     []model.Book `json:"books,omitempty"`
	CreatedAt time.Time    `json:"created_at"`
	UpdatedAt time.Time    `json:"updated_at"`
}
