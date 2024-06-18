package config

import "log-service/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
