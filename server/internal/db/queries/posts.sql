-- name: CreatePost :one
INSERT INTO posts (id, author_id, content, content_html, visibility, reply_to_post_id, original_post_id, language, meta)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: GetPost :one
SELECT * FROM posts WHERE id = $1;

-- name: UpdatePost :one
UPDATE posts SET content = $1, content_html = $2, visibility = $3, reply_to_post_id = $4, original_post_id = $5, language = $6, meta = $7, likes_count=$8, comments_count=$9, shares_count=$10, updated_at = $11 WHERE id = $12
RETURNING *;

-- name: DeletePost :one
DELETE FROM posts WHERE id = $1
RETURNING *;

