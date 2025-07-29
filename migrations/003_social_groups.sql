-- +goose Up
CREATE TABLE IF NOT EXISTS social_groups
(
    id          serial PRIMARY KEY,
    description text   NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE social_groups;