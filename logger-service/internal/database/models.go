// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package database

import (
	"time"
)

type Log struct {
	ID          int32
	ServiceName string
	LogData     string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
