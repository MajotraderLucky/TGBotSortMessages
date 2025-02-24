package main

import (
	"context"
	"log"
	"os"
	"strconv" // Добавлен strconv

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/session"
	"github.com/joho/godotenv"

	"tgmessenger/internal/auth"
)

// loadConfig загружает переменные окружения и возвращает API_ID, API_HASH
func loadConfig() (int, string, error) {
	if err := godotenv.Load(); err != nil {
		return 0, "", err
	}

	apiID, err := strconv.Atoi(os.Getenv("API_ID")) // strconv используется здесь
	if err != nil {
		return 0, "", err
	}

	apiHash := os.Getenv("API_HASH") // os используется здесь
	if apiHash == "" {
		return 0, "", err
	}

	return apiID, apiHash, nil
}

// initTelegramClient создает новый клиент Telegram
func initTelegramClient(apiID int, apiHash string) *telegram.Client {
	sessStorage := &session.FileStorage{Path: "session.json"}
	return telegram.NewClient(apiID, apiHash, telegram.Options{
		SessionStorage: sessStorage,
	})
}

func main() {
	apiID, apiHash, err := loadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	client := initTelegramClient(apiID, apiHash)

	ctx := context.Background()

	err = client.Run(ctx, func(ctx context.Context) error {
		log.Println("Бот запущен и подключен к Telegram")

		user, err := auth.AuthorizeClient(ctx, client)
		if err != nil {
			return err
		}

		log.Printf("✅ Успешно авторизован как %s", user.Username)
		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

