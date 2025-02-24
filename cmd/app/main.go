package main

import (
	"context"
	"log"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/config"
	"tgmessenger/internal/storage"
	"tgmessenger/internal/telegram"
)

func runBot(ctx context.Context, client *telegram.Client) error {
	log.Println("🚀 Бот запущен и подключен к Telegram")

	user, err := auth.AuthorizeClient(ctx, client.RawClient())
	if err != nil {
		return err
	}
	log.Printf("✅ Успешно авторизован как %s", user.Username)

	messages, err := telegram.GetUnreadMessages(ctx, client)
	if err != nil {
		log.Printf("⚠️ Ошибка получения сообщений: %v", err)
		return nil
	}

	if err := storage.SaveMessagesToJSON(messages); err != nil {
		log.Printf("⚠️ Ошибка сохранения сообщений: %v", err)
	}

	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	client := telegram.NewClient(cfg)
	ctx := context.Background()

	if err := client.RawClient().Run(ctx, func(ctx context.Context) error {
		return runBot(ctx, client)
	}); err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

