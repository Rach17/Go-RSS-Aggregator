-- name: CreateFeed :one
INSERT INTO feeds (title, url, description, language)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetFeedByID :one
SELECT * FROM feeds WHERE id = $1;

-- name: GetFeedByURL :one
SELECT * FROM feeds WHERE url = $1;

-- name: UpdateFeedLastFetchedAt :exec
UPDATE feeds
SET last_fetched_at = NOW()
WHERE id = $1;

