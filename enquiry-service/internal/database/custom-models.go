package database

import (
	"encoding/json"
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

type ScheduleJSON struct {
	ID            int32           `json:"id"`
	UserID        int32           `json:"user_id"`
	TaskType      string          `json:"task_type"`
	TaskDetails   json.RawMessage `json:"task_details"`
	ScheduledTime string          `json:"scheduled_time"`
	CreatedAt     time.Time       `json:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at"`
}

func ConvertTasksToJSON(tasks []Schedule) []ScheduleJSON {
	var tasksJSON []ScheduleJSON
	for _, task := range tasks {
		tasksJSON = append(tasksJSON, ScheduleJSON{
			ID:            task.ID,
			UserID:        task.UserID,
			TaskType:      task.TaskType,
			TaskDetails:   task.TaskDetails,
			ScheduledTime: task.ScheduledTime,
			CreatedAt:     task.CreatedAt,
			UpdatedAt:     task.UpdatedAt,
		})
	}
	return tasksJSON
}
