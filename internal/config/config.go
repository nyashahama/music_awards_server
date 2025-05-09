// config/config.go

package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

type DBConfig struct {
    Host, Port, User, Password, Name, SSLMode string
}

func mustGetenv(key string) (string, error) {
    v := os.Getenv(key)
    if v == "" {
        return "", fmt.Errorf("environment variable %q not set", key)
    }
    return v, nil
}

// LoadDBConfig reads .env (if present) and then requires each DB_ var.
func LoadDBConfig() (*DBConfig, error) {
    _ = godotenv.Load() // ignore error; we’ll catch missing vars below

    host, err := mustGetenv("DB_HOST");     if err != nil { return nil, err }
    port, err := mustGetenv("DB_PORT");     if err != nil { return nil, err }
    user, err := mustGetenv("DB_USER");     if err != nil { return nil, err }
    pass, err := mustGetenv("DB_PASSWORD"); if err != nil { return nil, err }
    name, err := mustGetenv("DB_NAME");     if err != nil { return nil, err }
    ssl,  err := mustGetenv("DB_SSLMODE");  if err != nil { return nil, err }

    return &DBConfig{
      Host: host, Port: port,
      User: user, Password: pass,
      Name: name, SSLMode: ssl,
    }, nil
}

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
    return db, nil
}

func (c *DBConfig) DatabaseURL() string {
    return fmt.Sprintf(
      "postgres://%s:%s@%s:%s/%s?sslmode=%s",
      c.User, c.Password, c.Host, c.Port, c.Name, c.SSLMode,
    )
}
