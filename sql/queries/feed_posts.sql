-- name: CreateFeedPost :exec
-- description: Create a new feed post
insert into feed_posts (feed_id, title, url, description,author, published_at)
values ($1, $2, $3, $4, $5, $6);

-- name: GetFeedPosts :many
select sqlc.embed(feeds), sqlc.embed(feed_posts) from feed_posts, feeds
where feed_posts.url = $1 and feed_posts.feed_id = feeds.id;

-- name: GetFeedPostByURL :one
select * from feed_posts where url = $1;