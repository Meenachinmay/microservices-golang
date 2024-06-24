-- +goose Up
-- Step 1: Add the columns without NOT NULL constraint
ALTER TABLE users
    ADD COLUMN available_timings VARCHAR(255),
    ADD COLUMN preferred_contact_method VARCHAR(255);

-- Step 2: Update existing rows to have default values
UPDATE users
SET available_timings = '09:00-18:00', -- Default available timings
    preferred_contact_method = 'email'; -- Default preferred contact method

-- Step 3: Alter the columns to add NOT NULL constraint
ALTER TABLE users
    ALTER COLUMN available_timings SET NOT NULL,
    ALTER COLUMN preferred_contact_method SET NOT NULL;

-- +goose Down
ALTER TABLE users
    DROP COLUMN available_timings,
    DROP COLUMN preferred_contact_method;
