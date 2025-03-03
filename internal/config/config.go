package config

import (
	"os"
	"strconv"
	"fmt"

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
		return nil, fmt.Errorf("ошибка загрузки .env файла: %w", err)
	}

	apiIDStr := os.Getenv("API_ID")
	if apiIDStr == "" {
		return nil, fmt.Errorf("переменная окружения API_ID отсутствует")
	}

	apiID, err := strconv.Atoi(apiIDStr)
	if err != nil {
		return nil, fmt.Errorf("неверный формат API_ID: %w", err)
	}

	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		return nil, fmt.Errorf("переменная окружения API_HASH отсутствует")
	}

	return &Config{
		APIID:   apiID,
		APIHash: apiHash,
	}, nil
}

