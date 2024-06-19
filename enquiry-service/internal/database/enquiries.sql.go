// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: enquiries.sql

package database

import (
	"context"
	"time"
)

const addNewEnquiry = `-- name: AddNewEnquiry :one
INSERT INTO enquiries (user_id, property_id, enquiry_date, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW(), NOW())
RETURNING id, user_id, property_id, enquiry_date, created_at, updated_at
`

type AddNewEnquiryParams struct {
	UserID     int32
	PropertyID int32
}

func (q *Queries) AddNewEnquiry(ctx context.Context, arg AddNewEnquiryParams) (Enquiry, error) {
	row := q.db.QueryRowContext(ctx, addNewEnquiry, arg.UserID, arg.PropertyID)
	var i Enquiry
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.PropertyID,
		&i.EnquiryDate,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getAllEnquiries = `-- name: GetAllEnquiries :many
SELECT id, user_id, property_id, enquiry_date, created_at, updated_at
FROM enquiries
`

func (q *Queries) GetAllEnquiries(ctx context.Context) ([]Enquiry, error) {
	rows, err := q.db.QueryContext(ctx, getAllEnquiries)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Enquiry
	for rows.Next() {
		var i Enquiry
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.PropertyID,
			&i.EnquiryDate,
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

const getPropertiesForEnquiry = `-- name: GetPropertiesForEnquiry :many
SELECT DISTINCT p.id, p.name, p.location, p.created_at, p.updated_at
FROM enquiries e
         JOIN properties p ON e.property_id = p.id
WHERE e.id = $1
`

func (q *Queries) GetPropertiesForEnquiry(ctx context.Context, id int32) ([]Property, error) {
	rows, err := q.db.QueryContext(ctx, getPropertiesForEnquiry, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Property
	for rows.Next() {
		var i Property
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Location,
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

const getUsersForEnquiry = `-- name: GetUsersForEnquiry :many
SELECT DISTINCT u.id, u.email, u.name, u.created_at, u.updated_at
FROM enquiries e
         JOIN users u ON e.user_id = u.id
WHERE e.id = $1
`

type GetUsersForEnquiryRow struct {
	ID        int32
	Email     string
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (q *Queries) GetUsersForEnquiry(ctx context.Context, id int32) ([]GetUsersForEnquiryRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsersForEnquiry, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersForEnquiryRow
	for rows.Next() {
		var i GetUsersForEnquiryRow
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Name,
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
