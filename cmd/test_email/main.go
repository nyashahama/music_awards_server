package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/nyashahama/music-awards/internal/config"
	"github.com/nyashahama/music-awards/internal/services"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load email config
	emailCfg := config.LoadEmailConfig()
	emailService := services.NewEmailService(emailCfg)

	// Test email
	log.Println("Sending test email...")
	err := emailService.SendWelcomeEmail("test@example.com", "Test User", 5)
	if err != nil {
		log.Printf("Failed to send test email: %v", err)

		// Check common issues
		if emailCfg.SMTPUsername == "" {
			log.Println("SMTP_USERNAME is empty")
		}
		if emailCfg.SMTPPassword == "" {
			log.Println("SMTP_PASSWORD is empty")
		}
		log.Printf("Using SMTP server: %s:%d", emailCfg.SMTPHost, emailCfg.SMTPPort)
	} else {
		log.Println("Test email sent successfully!")
	}
}
