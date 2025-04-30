package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

// DB is the globally accessible database connection.
var DB *sql.DB

// InitDB loads environment variables, opens a connection to Postgres,
// verifies it with Ping, and assigns the global DB.
func InitDB() *sql.DB {
	// Load .env if present (no fatal if missing)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	// Gather connection info from env vars
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	sslmode := os.Getenv("DB_SSLMODE")

	// Build the DSN
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, name, sslmode,
	)

	// Open the connection
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("Error opening DB: %v", err)
	}

	// Verify with Ping
	if err := db.Ping(); err != nil {
		log.Fatalf("Error pinging DB: %v", err)
	}

	log.Println("Connected to DB!")
	DB = db
	return db
}
