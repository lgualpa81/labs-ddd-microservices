package config

import "poc-product-svc/internal/infrastructure/utils"

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: utils.GetEnv("SERVER_PORT", "3000"),
		},
		Database: DatabaseConfig{
			Host:     utils.GetEnv("DB_HOST", "localhost"),
			Port:     utils.GetEnv("DB_PORT", "5432"),
			User:     utils.GetEnv("DB_USER", "postgres"),
			Password: utils.GetEnv("DB_PASSWORD", "postgres"),
			Name:     utils.GetEnv("DB_NAME", "product"),
			SSLMode:  utils.GetEnv("DB_SSLMODE", "disable"),
		},
	}
}
