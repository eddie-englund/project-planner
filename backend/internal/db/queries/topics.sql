-- name: CreateProjectTopic :one
INSERT INTO project_topics (id, project_id, index, title, color, image_url)
VALUES (gen_random_uuid(), $1, $2, $3, $4, $5)
RETURNING *;

-- name: GetProjectTopicByID :one
SELECT * FROM project_topics WHERE id = $1 AND project_id = $2;

-- name: ListProjectTopics :many
SELECT * FROM project_topics WHERE project_id = $1 ORDER BY index ASC;

-- name: UpdateProjectTopic :one
UPDATE project_topics SET title = $2, color = $3, image_url = $4, index = $5
WHERE id = $1 RETURNING *;

-- name: DeleteProjectTopic :exec
DELETE FROM project_topics WHERE id = $1;
