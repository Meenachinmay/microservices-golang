-- name: GetAllEnquiries :many
SELECT id, user_id, property_id, enquiry_date, created_at, updated_at
FROM enquiries;

-- name: GetUsersForEnquiry :many
SELECT DISTINCT u.id, u.email, u.name, u.created_at, u.updated_at
FROM enquiries e
         JOIN users u ON e.user_id = u.id
WHERE e.id = $1;

-- name: GetPropertiesForEnquiry :many
SELECT DISTINCT p.id, p.name, p.location, p.created_at, p.updated_at
FROM enquiries e
         JOIN properties p ON e.property_id = p.id
WHERE e.id = $1;

-- name: AddNewEnquiry :one
INSERT INTO enquiries (user_id, property_id, enquiry_date, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW(), NOW())
RETURNING *;