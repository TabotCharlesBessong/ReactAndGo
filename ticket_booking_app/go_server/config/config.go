package config

import (
	"path/filepath"

	"github.com/caarlos0/env"
	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type EnvConfig struct {
	ServerPort string `env:"SERVER_PORT,required"`
	DBHost     string `env:"DB_HOST,required"`
	DBName     string `env:"DB_NAME,required"`
	DBUser     string `env:"DB_USER,required"`
	DBPassword string `env:"DB_PASSWORD,required"`
	DBSSLMode  string `env:"DB_SSLMODE,required"`
}

func NewEnvConfig() *EnvConfig {
	// Try to load .env from the current directory first
	err := godotenv.Load()
	if err != nil {
		// If that fails, try to load from the project root
		err = godotenv.Load(filepath.Join("..", ".env"))
		if err != nil {
			log.Fatalf("Unable to load .env file: %v", err)
		}
	}

	config := &EnvConfig{}

	if err := env.Parse(config); err != nil {
		log.Fatalf("Unable to load variables from .env: %v", err)
	}

	return config
}
