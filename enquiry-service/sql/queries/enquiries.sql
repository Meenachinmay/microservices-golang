-- name: GetAllEnquiries :many
SELECT id, user_id, property_id, enquiry_date, created_at, updated_at
FROM enquiries;

-- name: GetAllEnquiriesForAPropertyByIdWithUsers :many
SELECT e.id, e.user_id, e.property_id, e.enquiry_date, e.created_at, e.updated_at, u.id, u.email, u.name
FROM enquiries e
         JOIN users u ON e.user_id = u.id
WHERE e.property_id = $1;

-- name: GetAllEnquiriesMadeByAUserByIdWithProperties :many
SELECT e.id, e.user_id, e.property_id, e.enquiry_date, e.created_at, e.updated_at, p.name, p.location
FROM enquiries e
         JOIN properties p ON e.property_id = p.id
WHERE e.user_id = $1;

-- name: CreateEnquiry :one
INSERT INTO enquiries (user_id, property_id, enquiry_date, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW(), NOW())
RETURNING *;