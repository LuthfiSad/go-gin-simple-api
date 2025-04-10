package dto

import (
	"time"

	"github.com/google/uuid"
)

type ChargeResponse struct {
	ID                uuid.UUID                `json:"id"`
	BookTransactionID uuid.UUID                `json:"book_transaction_id"`
	BookTransaction   *BookTransactionResponse `json:"book_transaction,omitempty"`
	DaysLate          int                      `json:"days_late"`
	DailyLateFee      float64                  `json:"daily_late_fee"`
	Total             float64                  `json:"total"`
	UserID            uuid.UUID                `json:"user_id"`
	User              *UserData                `json:"user,omitempty"`
	CreatedAt         time.Time                `json:"created_at"`
}

type ChargeCreateRequest struct {
	BookTransactionID uuid.UUID `json:"book_transaction_id" validate:"required"`
	DailyLateFee      float64   `json:"daily_late_fee" validate:"required,min=0"`
}

type ChargeUpdateRequest struct {
	DailyLateFee *float64 `json:"daily_late_fee,omitempty" validate:"omitempty,min=0"`
}
