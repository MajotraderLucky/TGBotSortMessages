package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config структура для хранения API_ID и API_HASH
type Config struct {
	APIID   int
	APIHash string
}

// LoadConfig загружает переменные окружения и возвращает Config
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, err
	}

	apiID, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		return nil, err
	}

	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		return nil, err
	}

	return &Config{
		APIID:   apiID,
		APIHash: apiHash,
	}, nil
}

