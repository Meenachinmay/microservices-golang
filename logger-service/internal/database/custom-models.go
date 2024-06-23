package database

import (
	"time"
)

type LogJSONResponse struct {
	ID          int32     `json:"id"`
	ServiceName string    `json:"service_name"`
	LogData     string    `json:"log_data"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type LogJSONPayload struct {
	ServiceName string `json:"service_name"`
	LogData     string `json:"log_data"`
}

func ConvertLogsToJSON(logs []Log) []LogJSONResponse {
	var logsJSON []LogJSONResponse

	for _, log := range logs {
		logsJSON = append(logsJSON, LogJSONResponse{
			ID:          log.ID,
			ServiceName: log.ServiceName,
			LogData:     log.LogData,
			CreatedAt:   log.CreatedAt,
			UpdatedAt:   log.UpdatedAt,
		})
	}
	return logsJSON
}
