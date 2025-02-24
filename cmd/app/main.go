package main

import (
	"log"
	"time"

	"tgmessenger/internal/app"
	"tgmessenger/internal/bot"
)

func main() {
	client, ctx, cancel := app.InitBot()
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

