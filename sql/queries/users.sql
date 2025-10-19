-- name: GetUserByEmail :one
SELECT * FROM users where email = $1;

-- name: CreateUser :one
INSERT INTO users (id, name, email, password,created_at,updated_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;