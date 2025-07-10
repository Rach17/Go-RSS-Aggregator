-- +goose Up
create table feeds (
   id              uuid primary key default gen_random_uuid(),
   created_at      timestamp with time zone default now() not null,
   updated_at      timestamp with time zone default null,
   title           text not null,
   url             text not null unique,
   description     text,
   language        text not null,
   last_fetched_at timestamp with time zone default now()
);

-- +goose Down
drop table feeds;