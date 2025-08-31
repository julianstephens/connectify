-- name: ListUserPostsOffset :many
SELECT * FROM posts WHERE author_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: CountUserPosts :one
SELECT COUNT(*) FROM posts WHERE author_id = $1;

-- name: ListUserPostsFirstPage :many
SELECT * FROM posts WHERE author_id = $1
ORDER BY created_at DESC
LIMIT $2;

-- name: ListUserPostsAfter :many
SELECT * FROM posts WHERE author_id = $1 AND created_at < $2
ORDER BY created_at DESC
LIMIT $3;