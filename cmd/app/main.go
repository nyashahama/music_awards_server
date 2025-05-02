package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/nyashahama/music-awards/internal/config"
	"github.com/nyashahama/music-awards/internal/db"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Load database configuration
	dbCfg, err := config.LoadDBConfig()
	if err != nil {
		log.Fatalf("Failed to load DB config: %v", err)
	}

	// 2. Initialize database connection
	sqlDB, err := config.InitDB(dbCfg)
	if err != nil {
		log.Fatalf("Failed to init DB: %v", err)
	}
	defer func() {
		if err := sqlDB.Close(); err != nil {
			log.Printf("Error closing database connection: %v", err)
		}
	}()

	// 3. Create GORM instance using existing connection
	gormDB, err := gorm.Open(
		postgres.New(postgres.Config{Conn: sqlDB}),
		&gorm.Config{},
	)
	if err != nil {
		log.Fatalf("Failed to create GORM instance: %v", err)
	}

	// 4. Perform database migrations
	if err := db.MigrateModels(gormDB,
		&models.User{},
		&models.Category{},
		&models.Nominee{},
		&models.Vote{},
	); err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	// 5. Set up graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
}