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
SELECT feeds.*, users.name as username
FROM feeds
INNER JOIN users
ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;

-- name: MarkFeedFetced :one
UPDATE feeds
SET updated_at = $1,last_fetched_at = $2
WHERE id = $3 RETURNING *;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1;

-- name: MarkAsFetched :one
UPDATE feeds
SET last_fetched_at = NOW(),
updated_at = NOW()
WHERE id = $1
RETURNING *;