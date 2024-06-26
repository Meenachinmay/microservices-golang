-- +goose Up
CREATE TABLE schedules (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    task_type VARCHAR(50) NOT NULL, -- can be 'call', 'sms' or 'email'
    task_details JSONB NOT NULL,
    scheduled_time TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id)
);

-- +goose Down
DROP TABLE IF EXISTS schedules;