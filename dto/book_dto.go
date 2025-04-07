package dto

import (
	"go-gin-simple-api/model"
	"time"

	"github.com/google/uuid"
)

// Book DTOs
type BookRes struct {
	ID          uuid.UUID    `json:"id"`
	Title       string       `json:"title"`
	Description string       `json:"description"`
	Cover       *model.Media `json:"cover,omitempty"`
	CoverURL    string       `json:"cover_url,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

type BookCreateReq struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Description string     `json:"description" validate:"max=1000"`
	CoverID     *uuid.UUID `json:"cover_id"`
}

type BookUpdateReq struct {
	Title       string     `json:"title" validate:"omitempty,max=255"`
	Description string     `json:"description" validate:"omitempty,max=1000"`
	CoverID     *uuid.UUID `json:"cover_id"`
}
