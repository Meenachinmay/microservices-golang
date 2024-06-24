-- name: AddNewEnquiryToUserById :one
UPDATE users
SET enquiry_count = enquiry_count + 1, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: GetUserByIdWithEnquiry :one
SELECT *
FROM users
WHERE id = $1;

-- name: CountEnquiriesForUserInLastWeek :one
SELECT COUNT(*)
FROM enquiries
WHERE user_id = $1
  AND enquiry_date >= $2;
