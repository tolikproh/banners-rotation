-- +goose Up
CREATE TABLE IF NOT EXISTS slots
(
    id          serial PRIMARY KEY,
    description text   NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE slots;