package main

import (
	"context"
	"log"
	"time"

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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	log.Println("🚀 Бот запущен! Проверка сообщений каждую минуту.")

	for {
		select {
		case <-ctx.Done():
			log.Println("⏹️ Завершаем работу бота...")
			return
		case <-ticker.C:
			if err := bot.Run(ctx, client); err != nil {
				log.Printf("⚠️ Ошибка работы бота: %v", err)
			}
		}
	}
}

