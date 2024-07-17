-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetPostsByUser :many
SELECT p.* FROM posts p
INNER JOIN feed_follows ff ON ff.feed_id = p.feed_id
WHERE ff.user_id = $1
ORDER BY published_at DESC, p.created_at DESC
LIMIT $2;

-- name: GetPostByFeedAndUrl :one
SELECT * FROM posts
WHERE feed_id = $1 AND url = $2
LIMIT 1;