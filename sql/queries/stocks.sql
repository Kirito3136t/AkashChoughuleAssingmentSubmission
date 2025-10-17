-- name: GetAllStocks :many
SELECT * 
FROM stocks;

-- name: GetStockById :one
SELECT * FROM stocks where ID = $1;

-- name: GetStockBySymbol :one
SELECT * FROM stocks where stock_symbol = $1;