package models

import "github.com/google/uuid"

type PortfolioRequest struct {
	UserID uuid.UUID `json:"user_id"`
}

type RecordPortfolioRequest struct {
	UserId   uuid.UUID `json:"user_id"`
	StockId  uuid.UUID `json:"stock_id"`
	Quantity string    `json:"quantity"`
	Type     string    `json:"type"`
}
