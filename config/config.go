package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	DatabasePath   string
	ServerAddress  string
	JWTSecret      string
	Environment    string

	// SMTP для уведомлений
	SMTPHost       string
	SMTPPort       string
	SMTPUser       string
	SMTPPass       string
	FromEmail      string
	NotifyAdminEmail string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Попытка прочитать .env файл (необязательно)
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			log.Println("Error reading config file:", err)
		}
	}

	cfg := &Config{
		DatabasePath:     viper.GetString("DATABASE_PATH"),
		ServerAddress:    viper.GetString("SERVER_ADDRESS"),
		JWTSecret:        viper.GetString("JWT_SECRET"),
		Environment:      viper.GetString("ENVIRONMENT"),
		SMTPHost:         viper.GetString("SMTP_HOST"),
		SMTPPort:         viper.GetString("SMTP_PORT"),
		SMTPUser:         viper.GetString("SMTP_USER"),
		SMTPPass:         viper.GetString("SMTP_PASS"),
		FromEmail:        viper.GetString("FROM_EMAIL"),
		NotifyAdminEmail: viper.GetString("NOTIFY_ADMIN_EMAIL"),
	}

	// Дефолтные значения
	if cfg.DatabasePath == "" {
		cfg.DatabasePath = "./data/gym_strongcode.db"
	}
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = ":8080"
	}
	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}
	if cfg.Environment == "" {
		cfg.Environment = "development"
	}

	return cfg
}