// Package config
package config

import (
	"os"
	"strconv"
)

type EmailConfig struct {
	SMTPHost     string
	SMTPPort     int
	SMTPUsername string
	SMTPPassword string
	FromEmail    string
	FromName     string
	UseTLS       bool
}

func LoadEmailConfig() *EmailConfig {
	port, _ := strconv.Atoi(getEnv("SMTP_PORT", "587"))
	useTLS, _ := strconv.ParseBool(getEnv("SMTP_USE_TLS", "true"))

	return &EmailConfig{
		SMTPHost:     getEnv("SMTP_HOST", "smtp.gmail.com"),
		SMTPPort:     port,
		SMTPUsername: getEnv("SMTP_USERNAME", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		FromEmail:    getEnv("FROM_EMAIL", "noreply@yourapp.com"),
		FromName:     getEnv("FROM_NAME", "Music Awards"),
		UseTLS:       useTLS,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
