package bot

import (
	"context"
	"log"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/storage"
	"tgmessenger/internal/telegram"
)

// Run запускает бота
func Run(ctx context.Context, client *telegram.Client) error {
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

