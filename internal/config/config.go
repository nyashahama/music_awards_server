package config // config/config.go

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
}

func Load() (*Config, error) {
	// Load from .env file or environment variables
}