-- name: CreateConnection :one
INSERT INTO connections (user_a, user_b, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserConnections :many
SELECT * FROM connections WHERE user_a = $1 OR user_b = $1;

-- name: UpdateConnection :one
UPDATE connections SET status = $1 WHERE user_a = $2 AND user_b = $3
RETURNING *;

-- name: DeleteConnection :one
DELETE FROM connections WHERE user_a = $1 AND user_b = $2
RETURNING *;