-- name: CreateStockTransaction :one
INSERT INTO stock_transactions(
    id, 
    user_id,
    stock_id,
    type,      
    quantity,           
    price,
    transaction_type,   
    created_at 
)
VALUES ($1,$2,$3,$4,$5,$6,$7,$8)
RETURNING *;

-- name: GetTodaysUserStock :many
SELECT * FROM stock_transactions
where user_id = $1 and DATE(created_at) = $2;