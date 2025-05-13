package app

import (
	"context"

	"tgmessenger/internal/config"
	"tgmessenger/internal/telegram"
)

// InitBot инициализирует Telegram-клиент с предоставленной конфигурацией
func InitBot(cfg *config.Config) (*telegram.Client, context.Context, context.CancelFunc) {
	client := telegram.NewClient(cfg)
	ctx, cancel := context.WithCancel(context.Background())

	return client, ctx, cancel
}
