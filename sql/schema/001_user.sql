-- +goose Up

create table users (
   id            uuid primary key default gen_random_uuid(),
   created_at    timestamp with time zone default now() not null,
   updated_at    timestamp with time zone default null,
   username      varchar(50) not null unique,
   password_hash varchar(255) not null
);

-- +goose Down
drop table users;