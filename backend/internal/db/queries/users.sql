-- backend/internal/db/queries/users.sql

-- name: CreateUser :one
INSERT INTO users (email)
VALUES ($1)
RETURNING id, email, created_at;

-- name: GetUserByID :one
SELECT id, email, created_at
FROM users
WHERE id = $1;