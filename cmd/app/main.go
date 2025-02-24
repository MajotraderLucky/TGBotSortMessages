package main

import (
	"context"
	"log"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/config"
	"tgmessenger/internal/telegram"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Создаём клиент Telegram
	client := telegram.NewClient(cfg)

	ctx := context.Background()

	// Запускаем клиента
	err = client.RawClient().Run(ctx, func(ctx context.Context) error {
		log.Println("Бот запущен и подключен к Telegram")

		user, err := auth.AuthorizeClient(ctx, client.RawClient())
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

