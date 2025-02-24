package app

import (
	"context"
	"log"

	"tgmessenger/internal/config"
	"tgmessenger/internal/telegram"
)

// InitBot инициализирует конфигурацию и Telegram-клиент
func InitBot() (*telegram.Client, context.Context, context.CancelFunc) {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	client := telegram.NewClient(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	return client, ctx, cancel
}

