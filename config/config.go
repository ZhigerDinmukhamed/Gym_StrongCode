// config/config.go
package config

import (
	"os"
)

type Config struct {
	DatabasePath  string
	ServerAddress string
	JWTSecret     []byte
	Environment   string
}

// Load загружает конфигурацию из переменных окружения
func Load() *Config {
	cfg := &Config{
		DatabasePath:  getEnv("DATABASE_PATH", "gym_strongcode.db"),
		ServerAddress: getEnv("SERVER_ADDRESS", ":8080"),
		JWTSecret:     []byte(getEnv("JWT_SECRET", "strongcode-secret-change-in-production")),
		Environment:   getEnv("ENVIRONMENT", "development"),
	}

	return cfg
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
