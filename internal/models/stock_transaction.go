package models

import "github.com/google/uuid"

type TransactionRequest struct {
	UserID   uuid.UUID `json:"user_id"`
	StockId  uuid.UUID `json:"stock_id"`
	Type     string    `json:"type"`
	Quantity string    `json:"quantity"`
	Price    string    `json:"price"`
}

type RequestBodyTransaction struct {
	Type     string  `json:"type"`
	Quantity float64 `json:"quantity"`
}

type RewardRequest struct {
	UserID      uuid.UUID `json:"user_id"`
	StockSymbol string    `json:"stock_symbol"`
	Quantity    string    `json:"quantity"`
}
