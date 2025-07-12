-- +goose Up
create table feed_posts (
    id              uuid primary key default gen_random_uuid(),
    created_at      timestamp with time zone default now() not null,
    updated_at      timestamp with time zone default null,
    feed_id         uuid not null references feeds(id) on delete cascade,
    title           text not null,
    url             text not null unique,
    description     text,
    published_at    timestamp with time zone not null,
    author          text
);

-- +goose Down
drop table feed_posts;