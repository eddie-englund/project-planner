-- name: CreateProject :one
INSERT INTO projects (id, title, color, image_url)
VALUES (gen_random_uuid(), $1, $2, $3)
RETURNING *;

-- name: GetProjectByID :one
SELECT * FROM projects WHERE id = $1;

-- name: ListProjects :many
SELECT * FROM projects ORDER BY created_at DESC;

-- name: UpdateProject :one
UPDATE projects SET title = $2, color = $3, image_url = $4
WHERE id = $1 RETURNING *;

-- name: DeleteProject :exec
DELETE FROM projects WHERE id = $1;
