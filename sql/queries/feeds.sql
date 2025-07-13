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
WHERE url = $1;

-- name: GetAllFeeds :many
SELECT * FROM feeds;


-- name: FollowFeed :exec
INSERT INTO feed_follow (user_id, feed_id)
VALUES ($1, $2)
ON CONFLICT (user_id, feed_id) DO NOTHING;

-- name: GetLastFetchedFeeds :many
SELECT * FROM feeds
ORDER BY last_fetched_at DESC
LIMIT $1;
