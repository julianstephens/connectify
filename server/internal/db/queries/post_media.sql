-- name: CreatePostMedia :one
INSERT INTO post_media (id, post_id, url, media_type, width, height, size_bytes, meta, sort_index)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetPostMedia :many
SELECT * FROM post_media WHERE post_id = $1;

-- name: DeletePostMedia :one
DELETE FROM post_media WHERE id = $1
RETURNING *;
