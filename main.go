// main.go
package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nyashahama/music-awards/internal/db"
	
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database connection
	pgConfig := &db.Config{
		Host:     cfg.DBHost,
		Port:     cfg.DBPort,
		User:     cfg.DBUser,
		Password: cfg.DBPassword,
		DBName:   cfg.DBName,
		SSLMode:  "disable",
	}

	// Create GORM connection
	gormDB, err := db.NewGormConnection(pgConfig)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}
	defer db.CloseConnection()

	// Run migrations
	err = db.MigrateModels(gormDB, &models.User{}, &models.Category{}, &models.Nominee{}, &models.Vote{})
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// Start server
	// ...

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}