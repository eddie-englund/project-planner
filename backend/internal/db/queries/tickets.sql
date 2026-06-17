-- name: CreateTicket :one
INSERT INTO topic_tickets (topic_id, status_id, title, body, urls)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: ListTicketsByTopic :many
SELECT * FROM topic_tickets WHERE topic_id = $1 ORDER BY created_at ASC;

-- name: GetTicketByID :one
SELECT * FROM topic_tickets WHERE id = $1 AND topic_id = $2;

-- name: UpdateTicket :one
UPDATE topic_tickets SET title = $2, body = $3, urls = $4, status_id = $5
WHERE id = $1 RETURNING *;

-- name: DeleteTicket :exec
DELETE FROM topic_tickets WHERE id = $1;
