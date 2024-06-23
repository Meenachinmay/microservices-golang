-- name: InsertLog :one
INSERT INTO logs (service_name, log_data)
VALUES ($1, $2)
RETURNING id, service_name, log_data, created_at, updated_at;

-- name: GetAllLogs :many
SELECT id, service_name, log_data, created_at, updated_at
FROM logs;