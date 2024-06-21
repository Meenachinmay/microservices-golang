package database

import (
	"time"
)

type PropertyJSON struct {
	ID         int32     `json:"id"`
	Name       string    `json:"name"`
	Location   string    `json:"location"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	FudousanID int32     `json:"fudousan_id"`
}

func ConvertPropertiesToJSON(properties []Property) []PropertyJSON {
	var propertiesJSON []PropertyJSON

	for _, property := range properties {
		propertiesJSON = append(propertiesJSON, PropertyJSON{
			ID:         property.ID,
			Name:       property.Name,
			Location:   property.Location,
			CreatedAt:  property.CreatedAt,
			UpdatedAt:  property.UpdatedAt,
			FudousanID: property.FudousanID,
		})
	}
	return propertiesJSON
}
