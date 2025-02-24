package main

import (
	"context"
	"log"
	"time"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/config"
	"tgmessenger/internal/storage"
	"tgmessenger/internal/telegram"
)

func runBot(ctx context.Context, client *telegram.Client) error {
	log.Println("🚀 Проверяем новые сообщения...")

	user, err := auth.AuthorizeClient(ctx, client.RawClient())
	if err != nil {
		return err
	}
	log.Printf("✅ Авторизован как %s", user.Username)

	messages, err := telegram.GetUnreadMessages(ctx, client)
	if err != nil {
		log.Printf("⚠️ Ошибка получения сообщений: %v", err)
		return nil
	}

	if err := storage.SaveMessagesToJSON(messages); err != nil {
		log.Printf("⚠️ Ошибка сохранения JSON: %v", err)
	}

	if err := storage.SaveMessagesToMarkdown(messages); err != nil {
		log.Printf("⚠️ Ошибка сохранения Markdown: %v", err)
	}

	log.Println("✅ Обновление завершено.")
	return nil
}

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
			if err := runBot(ctx, client); err != nil {
				log.Printf("⚠️ Ошибка работы бота: %v", err)
			}
		}
	}
}

