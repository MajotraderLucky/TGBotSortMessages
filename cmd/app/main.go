package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/config"
	"tgmessenger/internal/telegram"
)

// Сохраняем сообщения в JSON
func saveMessagesToJSON(messages []string) error {
	dir := "messages"
	filePath := filepath.Join(dir, "unread.json")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	log.Printf("📂 Сообщения сохранены в %s", filePath)
	return nil
}

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	client := telegram.NewClient(cfg)
	ctx := context.Background()

	err = client.RawClient().Run(ctx, func(ctx context.Context) error {
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

		if err := saveMessagesToJSON(messages); err != nil {
			log.Printf("⚠️ Ошибка сохранения сообщений: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

