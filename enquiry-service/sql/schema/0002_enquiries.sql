-- +goose Up
CREATE TABLE enquiries
(
    id           SERIAL PRIMARY KEY,
    user_id      INT       NOT NULL,
    property_id  INT       NOT NULL,
    enquiry_date TIMESTAMP NOT NULL DEFAULT NOW(),
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    FOREIGN KEY (user_id) REFERENCES users (id),
    FOREIGN KEY (property_id) REFERENCES properties (id)
);

-- +goose Down
DROP TABLE IF EXISTS enquiries;