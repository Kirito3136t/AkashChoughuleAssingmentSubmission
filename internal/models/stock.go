package models

type StockResponseObject struct {
	ID          string `json:"id"`
	StockSymbol string `json:"stock_symbol"`
	Valuation   string `json:"valuation"`
}
