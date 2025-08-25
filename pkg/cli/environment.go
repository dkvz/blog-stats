package cli

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DbPath string
}

func ConfigFromDotEnv() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	dbPath := os.Getenv("DB_PATH")
	return &Config{
		DbPath: dbPath,
	}, nil
}
