-- name: CreateFeedPost :exec
-- description: Create a new feed post
insert into feed_post (feed_id, title, url, description, published_at, author)
values ($1, $2, $3, $4, $5, $6);

-- name: GetFeedPosts :many
select * from feed_post, feeds 
where feed_post.url = $1 and feed_post.feed_id = feed.id;
