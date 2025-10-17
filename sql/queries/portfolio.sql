-- name: GetUserStock :one
SELECT * FROM portfolio where user_id = $1 and stock_id = $2;

-- name: RecordUserPortfolio :one
INSERT INTO portfolio (
    id,
    user_id,
    stock_id,
    total_quantity,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6
)
ON CONFLICT (user_id, stock_id)
DO UPDATE SET
    total_quantity = portfolio.total_quantity + EXCLUDED.total_quantity,
    updated_at = EXCLUDED.updated_at
RETURNING *;

-- name: GetPortfolioByUserId :many
Select * from portfolio where user_id = $1;
