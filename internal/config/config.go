// Package config
package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// DBConfig holds the settings for your Postgres connection.
type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// LoadDBConfig reads environment variables (and .env) into a DBConfig.
func LoadDBConfig() (*DBConfig, error) {
	// Load .env, but don’t fatal if missing
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	cfg := &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	// You could validate here that none of the strings are empty, if you like.
	return cfg, nil
}

// InitDB takes a DBConfig and opens/pings a *sql.DB.
func InitDB(cfg *DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.Name, cfg.SSLMode,
	)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("pinging DB: %w", err)
	}
	log.Println("Connected to DB!")
	return db, nil
}

// DatabaseURL builds a postgres:// URL string for golang-migrate.
func (c *DBConfig) DatabaseURL() string {
	// e.g. "postgres://user:pass@host:port/dbname?sslmode=disable"
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, c.Password,
		c.Host, c.Port,
		c.Name, c.SSLMode,
	)
}
