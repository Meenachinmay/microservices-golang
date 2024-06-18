package data

type LogEntry struct {
	ServiceName string `json:service_name`
	Data        string `json:data`
}
