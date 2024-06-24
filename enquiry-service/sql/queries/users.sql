-- name: AddNewEnquiryToUserById :one
UPDATE users
SET enquiry_count = enquiry_count + 1, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetUserById :one
SELECT *
FROM users
WHERE id = $1;

-- name: CountEnquiriesForUserInLastWeek :one
SELECT COUNT(*)
FROM enquiries
WHERE user_id = $1
  AND enquiry_date >= $2;

-- name: CreateUser :one
INSERT INTO users (email, name, available_timings, preferred_contact_method)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetUserByEmail :one
SELECT *
FROM users
WHERE email = $1;