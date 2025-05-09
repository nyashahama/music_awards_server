package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nyashahama/music-awards/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	dbPool     *pgxpool.Pool
	gormDB     *gorm.DB
	dbOnce     sync.Once
)

// Database configuration structure
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// NewPostgresPool creates a singleton PostgreSQL connection pool
func NewPostgresPool(config *Config) (*pgxpool.Pool, error) {
	var initErr error
	dbOnce.Do(func() {
		connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
			config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode)

		poolConfig, err := pgxpool.ParseConfig(connStr)
		if err != nil {
			initErr = err
			return
		}

		// Connection pool settings
		poolConfig.MaxConns = 50
		poolConfig.MinConns = 10
		poolConfig.MaxConnLifetime = time.Hour
		poolConfig.MaxConnIdleTime = time.Minute * 30
		poolConfig.HealthCheckPeriod = time.Minute * 5

		dbPool, err = pgxpool.NewWithConfig(context.Background(), poolConfig)
		if err != nil {
			initErr = err
			return
		}

		// Verify connection
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = dbPool.Ping(ctx); err != nil {
			initErr = err
			return
		}

		log.Println("Successfully connected to PostgreSQL!")
	})

	return dbPool, initErr
}

// NewGormConnection creates a GORM database connection
func NewGormConnection(config *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		config.Host, config.User, config.Password, config.DBName, config.Port, config.SSLMode)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		PrepareStmt:            true,
		SkipDefaultTransaction: true,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Set connection pool parameters
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	

	log.Println("GORM database connection established")
	return db, nil
}

// CloseConnection closes database connections
func CloseConnection() {
	if dbPool != nil {
		dbPool.Close()
		log.Println("Closed PostgreSQL connection pool")
	}

	if gormDB != nil {
		sqlDB, _ := gormDB.DB()
		sqlDB.Close()
		log.Println("Closed GORM database connection")
	}
}

// MigrateModels runs database migrations
func MigrateModels(db *gorm.DB) error {
   

    // Create tables with constraints directly
  
	err := db.Set("gorm:table_options", "WITHOUT OIDS").AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Nominee{},
		&models.NomineeCategory{},
		&models.Vote{},
	)
    if err != nil {
        return fmt.Errorf("failed to create tables: %w", err)
    }

    // Add foreign key constraints using ALTER TABLE
    foreignKeys := []struct {
        name    string
        sql     string
    }{
        // Your existing foreign key definitions
    }

    for _, fk := range foreignKeys {
        if err := db.Exec(fk.sql).Error; err != nil {
            // Handle error or log if constraint already exists
            log.Printf("Constraint %s might already exist: %v", fk.name, err)
        }
    }

    return nil
}
// HealthCheck verifies database connectivity
func HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return dbPool.Ping(ctx)
}