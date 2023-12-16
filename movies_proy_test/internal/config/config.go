package config

import (
	"fmt"
	"os"
)

// Config represents the application configuration.
// It includes settings which are crucial for the operation of the application,
// such as database connection information, JWT secrets, and server details.
type Config struct {
	DatabaseURL   string // DatabaseURL is the connection string for the PostgreSQL database
	JWTSecret     string // JWTSecret is the secret key used for signing JWTs
	ServerAddress string // ServerAddress is the address where the server listens
}

// LoadConfig loads application configuration from environment variables.
// It optionally reads from a .env file (useful in development environments).
// Returns a Config struct and an error if the loading fails.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
	}

	// Verifica que las configuraciones necesarias estén presentes
	if cfg.DatabaseURL == "" || cfg.JWTSecret == "" || cfg.ServerAddress == "" {
		return nil, fmt.Errorf("missing required environment variables")
	}

	return cfg, nil
}
