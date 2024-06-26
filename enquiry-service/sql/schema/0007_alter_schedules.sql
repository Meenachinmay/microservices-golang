-- +goose Up
ALTER TABLE schedules
    ALTER COLUMN scheduled_time TYPE VARCHAR(255),
    ALTER COLUMN scheduled_time SET NOT NULL;

-- +goose Down
ALTER TABLE schedules
    ALTER COLUMN scheduled_time TYPE TIMESTAMP,
    ALTER COLUMN scheduled_time DROP NOT NULL;