package models

import (
	"time"

	"github.com/google/uuid"
)

type PortfolioRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type RecordPortfolioRequest struct {
	UserId   uuid.UUID `json:"user_id"`
	StockId  uuid.UUID `json:"stock_id"`
	Quantity string    `json:"quantity"`
	Type     string    `json:"type"`
}

type ResponseObjectPortfolio struct {
	ID               uuid.UUID `json:"id"`
	Quantity         string    `json:"quantity"`
	CurrentValuation float64   `json:"valuation"`
	StockSymbol      string    `json:"stock_symbol"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
