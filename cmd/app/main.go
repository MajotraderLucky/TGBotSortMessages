package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/session"
	"github.com/joho/godotenv"

	"tgmessenger/internal/auth"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	apiID, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		log.Fatal("Ошибка конвертации API_ID в int:", err)
	}
	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		log.Fatal("API_HASH не указан в .env")
	}

	sessStorage := &session.FileStorage{Path: "session.json"}
	client := telegram.NewClient(apiID, apiHash, telegram.Options{
		SessionStorage: sessStorage,
	})

	ctx := context.Background()

	err = client.Run(ctx, func(ctx context.Context) error {
		log.Println("Бот запущен и подключен к Telegram")

		user, err := auth.AuthorizeClient(ctx, client) // Используем новую функцию
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

