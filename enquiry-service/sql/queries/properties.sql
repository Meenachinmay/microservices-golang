-- name: GetAPropertyDetailsById :one
SELECT name, location
FROM properties
WHERE id = $1;

-- name: FetchAllProperties :many
SELECT *
FROM properties;

-- name: GetAllPropertiesForAFudousan :many
SELECT *
FROM properties
WHERE fudousan_id = $1;