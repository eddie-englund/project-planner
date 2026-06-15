-- backend/internal/db/queries/users.sql

-- name: CreateUser :one
INSERT INTO users (email)
VALUES ($1)
RETURNING id, email, email_verified, created_at;

-- name: GetUserByID :one
SELECT id, email, email_verified, created_at
FROM users
WHERE id = $1;

-- name: UserExistsByEmail :one
SELECT EXISTS (
  SELECT 1
  FROM users
  WHERE email = $1
);