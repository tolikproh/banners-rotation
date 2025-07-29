-- +goose Up
INSERT INTO banners (description)
VALUES ('Youth banner'),
       ('Sport banner'),
       ('Health banner');
INSERT INTO slots (description)
VALUES ('Top slot'),
       ('Middle slot'),
       ('Bottom banner');
INSERT INTO social_groups (description)
VALUES ('Молодежь'),
       ('Люди среднего возраста'),
       ('Пожилые');

-- +goose Down
TRUNCATE TABLE banners;
TRUNCATE TABLE slots;
TRUNCATE TABLE social_groups;
