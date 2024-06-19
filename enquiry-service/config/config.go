package config

import "enquiry-service/internal/database"

type ApiConfig struct {
	DB *database.Queries
}
