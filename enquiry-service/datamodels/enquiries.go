package datamodels

import "time"

type CreateEnquiryType struct {
	UserID     int32 `json:"user_id"`
	PropertyID int32 `json:"property_id"`
}

type LoadEnquiryType struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id"`
	PropertyID  int32     `json:"property_id"`
	EnquiryDate time.Time `json:"enquiry_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
