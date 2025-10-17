-- name: CreateStockTransaction :one
INSERT INTO stock_transactions(
    id, 
    user_id,
    stock_id,
    type,      
    quantity,           
    price,   
    created_at 
)
VALUES ($1,$2,$3,$4,$5,$6,$7)
RETURNING *;