-- name: CreateUser :one
INSERT INTO users (username, password) VALUES ($1, $2) RETURNING *;

-- name: GetUserByUsername :one
SELECT * FROM users WHERE username = $1;

-- name: UpdateUserLoginTime :exec 
UPDATE users SET last_login = CURRENT_TIMESTAMP WHERE username = $1; 
