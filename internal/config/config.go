package config

import (
	"database/sql"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

var DB *sql.DB

func InitDB() *sql.DB {
	_ = godotenv.Load()
	dsn := "postgres://postgres:Gyver@localhost:5432/ecommerce_db?sslmode=disable"
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal("Error opening DB:", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal("Error pinging DB:", err)
	}
	log.Println("Connected to DB")
	DB = db
	return db
}
