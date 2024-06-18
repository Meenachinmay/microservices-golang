// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: logs.sql

package database

import (
	"context"
)

const getLogs = `-- name: GetLogs :many
SELECT id, service_name, log_data, created_at, updated_at
FROM logs
`

func (q *Queries) GetLogs(ctx context.Context) ([]Log, error) {
	rows, err := q.db.QueryContext(ctx, getLogs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Log
	for rows.Next() {
		var i Log
		if err := rows.Scan(
			&i.ID,
			&i.ServiceName,
			&i.LogData,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertLog = `-- name: InsertLog :one
INSERT INTO logs (service_name, log_data)
VALUES ($1, $2)
RETURNING id, service_name, log_data, created_at, updated_at
`

type InsertLogParams struct {
	ServiceName string
	LogData     string
}

func (q *Queries) InsertLog(ctx context.Context, arg InsertLogParams) (Log, error) {
	row := q.db.QueryRowContext(ctx, insertLog, arg.ServiceName, arg.LogData)
	var i Log
	err := row.Scan(
		&i.ID,
		&i.ServiceName,
		&i.LogData,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
