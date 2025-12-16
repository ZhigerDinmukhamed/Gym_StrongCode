package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	DatabasePath  string
	ServerAddress string
	JWTSecret     []byte
	Environment   string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv() // читает из окружения (важно для Docker!)
	_ = viper.ReadInConfig()

	cfg := &Config{
		DatabasePath:  viper.GetString("DATABASE_PATH"),
		ServerAddress: viper.GetString("SERVER_ADDRESS"),
		JWTSecret:     []byte(viper.GetString("JWT_SECRET")),
		Environment:   viper.GetString("ENVIRONMENT"),
	}

	// Дефолты (на случай запуска без .env)
	if cfg.DatabasePath == "" {
		cfg.DatabasePath = "gym_strongcode.db"
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8080"
	}
	if len(cfg.JWTSecret) == 0 {
		cfg.JWTSecret = []byte("strongcode-secret-change-in-production")
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}

	return cfg
}