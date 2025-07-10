-- +goose Up
CREATE TABLE feeds (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_at timestamp with time zone DEFAULT null,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    url text NOT NULL UNIQUE,
    title text NOT NULL,
    description text
);
-- +goose Down