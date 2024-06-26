-- name: CreateSchedule :one
INSERT INTO schedules (user_id, task_type, task_details, scheduled_time, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING id, user_id, task_type, task_details, scheduled_time, created_at, updated_at;

-- name: GetDueSchedules :many
SELECT id, user_id, task_type, task_details, scheduled_time, created_at, updated_at
FROM schedules
ORDER BY scheduled_time ASC;

-- name: DeleteSchedule :exec
DELETE FROM schedules
WHERE id = $1;
