-- +goose Up
ALTER TABLE users
    ADD COLUMN enquiry_count INT DEFAULT 0;

-- +goose Down
ALTER TABLE users
    DROP COLUMN enquiry_count;