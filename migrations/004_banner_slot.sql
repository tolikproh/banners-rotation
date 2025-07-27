-- +goose Up
CREATE TABLE IF NOT EXISTS banner_slot
(
    banner_id serial NOT NULL REFERENCES banners (id) MATCH FULL ON UPDATE CASCADE ON DELETE CASCADE,
    slot_id   serial NOT NULL REFERENCES slots (id) MATCH FULL ON UPDATE CASCADE ON DELETE CASCADE,
    PRIMARY KEY (banner_id, slot_id)
);

-- +goose Down
DROP TABLE banner_slot;