-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    name text NOT NULL,
    age bigserial NOT NULL
);

-- +goose Down
DROP TABLE users;
