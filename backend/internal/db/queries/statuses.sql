-- name: CreateStatus :one
INSERT INTO project_statuses (project_id, name, color, position, is_terminal)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListProjectStatuses :many
SELECT * FROM project_statuses WHERE project_id = $1 ORDER BY position ASC;

-- name: SeedDefaultStatuses :exec
INSERT INTO project_statuses (project_id, name, color, position, is_terminal)
VALUES
  ($1, 'Open',   '#6B8E7C', 0, false),
  ($1, 'Closed', '#8E7C6B', 1, true);
