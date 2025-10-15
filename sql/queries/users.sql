-- name: GetUserByEmail :one
Select * from users where email = $1;

-- name: CreateUser :one
INSERT INTO users (id, name, email, isReferral, referralUserId) 
VALUES ($1,$2,$3,$4,$5)
RETURNING *;;