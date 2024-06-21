-- +goose Up
ALTER TABLE properties
    ADD COLUMN fudousan_id INT NOT NULL DEFAULT 0;

-- +goose Down
ALTER TABLE properties
    DROP COLUMN fudousan_id;