-- name: GetAPropertyDetailsById :one
SELECT name, location
FROM properties
WHERE id = $1;