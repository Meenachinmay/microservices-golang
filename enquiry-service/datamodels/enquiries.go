package datamodels

import "time"

type LoadEnquiryType struct {
	ID          int32     `json:"id"`
	UserID      int32     `json:"user_id"`
	PropertyID  int32     `json:"property_id"`
	EnquiryDate time.Time `json:"enquiry_date"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type EnquiryPayload struct {
	UserID     string `json:"user_id"`
	PropertyID string `json:"property_id"`
	Name       string `json:"name"`
	Location   string `json:"location"`
}
