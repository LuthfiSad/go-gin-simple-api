package dto

import (
	"time"

	"github.com/google/uuid"
)

type CustomerResponse struct {
	ID        uuid.UUID `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerWithTransactionsResponse struct {
	ID               uuid.UUID                 `json:"id"`
	Code             string                    `json:"code"`
	Name             string                    `json:"name"`
	CreatedAt        time.Time                 `json:"created_at"`
	UpdatedAt        time.Time                 `json:"updated_at"`
	BookTransactions []BookTransactionResponse `json:"book_transactions,omitempty"`
}

type CustomerCreateRequest struct {
	Code string `json:"code" validate:"required,min=3,max=50"`
	Name string `json:"name" validate:"required,min=3,max=255"`
}

type CustomerUpdateRequest struct {
	Code *string `json:"code,omitempty" validate:"omitempty,min=3,max=50"`
	Name *string `json:"name,omitempty" validate:"omitempty,min=3,max=255"`
}
