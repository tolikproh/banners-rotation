-- +goose Up
CREATE TABLE IF NOT EXISTS banner_clicks
(
    id              serial    PRIMARY KEY,
    banner_id       serial    NOT NULL REFERENCES banners (id) MATCH FULL ON UPDATE CASCADE ON DELETE CASCADE,
    slot_id         serial    NOT NULL REFERENCES slots (id) MATCH FULL ON UPDATE CASCADE ON DELETE CASCADE,
    social_group_id serial    NOT NULL REFERENCES social_groups (id) MATCH FULL ON UPDATE CASCADE ON DELETE CASCADE,
    date            timestamp NOT NULL DEFAULT current_timestamp
);

-- +goose Down
DROP TABLE banner_clicks;