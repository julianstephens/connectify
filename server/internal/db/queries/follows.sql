-- name: CreateFollow :one
INSERT INTO follows (follower_id, followee_id, status)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateFollow :one
UPDATE follows SET status = $1 WHERE follower_id = $2 AND followee_id = $3
RETURNING *;

-- name: GetUserFollows :many
SELECT * FROM follows WHERE follower_id = $1;

-- name: DeleteFollow :one
DELETE FROM follows WHERE follower_id = $1 AND followee_id = $2
RETURNING *;
