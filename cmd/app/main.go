package main

import (
	"context"
	"log"

	"tgmessenger/internal/bot"
	"tgmessenger/internal/config"
	"tgmessenger/internal/telegram"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	client := telegram.NewClient(cfg)
	ctx := context.Background()

	if err := client.RawClient().Run(ctx, func(ctx context.Context) error {
		return bot.Run(ctx, client) // 💡 Используем bot.Run вместо runBot
	}); err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

