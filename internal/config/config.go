package config

import (
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config структура для хранения конфигурации приложения
type Config struct {
	APIID             int
	APIHash           string
	MessagesPerDialog int      // Количество сообщений для получения из каждого диалога
	UpdateInterval    int      // Интервал обновления в секундах
	OutputFormats     []string // Форматы вывода (json, markdown)
	Debug             bool     // Режим отладки
}

// LoadConfig загружает переменные окружения и конфигурацию из командной строки
func LoadConfig() (*Config, error) {
	// Загружаем .env файл если существует
	godotenv.Load() // Игнорируем ошибку, т.к. файл может и не существовать

	// Настройки по умолчанию
	config := &Config{
		MessagesPerDialog: 5,
		UpdateInterval:    60,
		OutputFormats:     []string{"json", "markdown"},
		Debug:             false,
	}

	// Определяем флаги командной строки
	messagesPerDialog := flag.Int("msgs", config.MessagesPerDialog, "Количество сообщений для получения из каждого диалога")
	updateInterval := flag.Int("interval", config.UpdateInterval, "Интервал обновления в секундах")
	outputFormats := flag.String("formats", "json,markdown", "Форматы вывода (через запятую): json,markdown")
	debug := flag.Bool("debug", config.Debug, "Включить подробный вывод отладочной информации")

	// Парсим аргументы
	flag.Parse()

	// Обновляем конфигурацию
	config.MessagesPerDialog = *messagesPerDialog
	config.UpdateInterval = *updateInterval
	config.Debug = *debug

	// Парсим форматы вывода
	if *outputFormats != "" {
		config.OutputFormats = []string{}
		for _, format := range splitCSV(*outputFormats) {
			if isValidFormat(format) {
				config.OutputFormats = append(config.OutputFormats, format)
			}
		}
	}

	// Загружаем API_ID и API_HASH из переменных окружения
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

	config.APIID = apiID
	config.APIHash = apiHash

	return config, nil
}

// Вспомогательная функция для разделения значений, указанных через запятую
func splitCSV(s string) []string {
	var result []string
	current := ""

	for _, c := range s {
		if c == ',' {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(c)
		}
	}

	if current != "" {
		result = append(result, current)
	}

	return result
}

// Проверяет, является ли формат допустимым
func isValidFormat(format string) bool {
	validFormats := map[string]bool{
		"json":       true,
		"markdown":   true,
		"structured": true,
		"byuser":     true,
	}

	return validFormats[format]
}
