-- +goose Up
create table feed_follow (
    id              uuid primary key default gen_random_uuid(),
    created_at      timestamp with time zone default now() not null,
    updated_at      timestamp with time zone default null,
    user_id         uuid not null references users(id) on delete cascade,
    feed_id         uuid not null references feeds(id) on delete cascade,
    unique (user_id, feed_id)
);

-- +goose Down
drop table feed_follow;