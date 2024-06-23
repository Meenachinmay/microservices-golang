-- +goose Up
CREATE TABLE logs (
                      id SERIAL PRIMARY KEY,
                      service_name VARCHAR(255) NOT NULL,
                      log_data TEXT NOT NULL,
                      created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
                      updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE IF EXISTS logs;
