-- name: CreateFeed :one
INSERT INTO feeds (id,name,created_at,updated_at,url,user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
) 
RETURNING *;

-- name: ListFeeds :many
SELECT *, users.name as username
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;