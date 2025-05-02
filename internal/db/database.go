package db

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
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

	// Enable UUID extension
	err = db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\"").Error
	if err != nil {
		return nil, fmt.Errorf("failed to create uuid extension: %w", err)
	}

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
func MigrateModels(db *gorm.DB, models ...interface{}) error {
    // First create UUID extension
    if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp"`).Error; err != nil {
        return fmt.Errorf("failed to create uuid extension: %w", err)
    }

    // Disable foreign key constraints during migration
    db.Exec("SET session_replication_role = replica;")
    defer db.Exec("SET session_replication_role = DEFAULT;")

    // Migrate tables with explicit settings
    if err := db.AutoMigrate(models...); err != nil {
        return fmt.Errorf("failed to migrate models: %w", err)
    }

    // Add composite foreign key for votes table
    if err := db.Exec(`
        ALTER TABLE votes
        ADD CONSTRAINT fk_vote_nominee_category
        FOREIGN KEY (nominee_id, category_id) 
        REFERENCES nominee_categories(nominee_id, category_id)
        ON DELETE CASCADE
    `).Error; err != nil {
        return fmt.Errorf("failed to create composite foreign key: %w", err)
    }

    return nil
}

// HealthCheck verifies database connectivity
func HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return dbPool.Ping(ctx)
}