-- +goose Up
CREATE TABLE IF NOT EXISTS banners
(
    id          serial PRIMARY KEY,
    description text   NOT NULL DEFAULT ''
);

-- +goose Down
DROP TABLE banners;