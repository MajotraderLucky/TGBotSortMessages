package bot

import (
	"context"
	"log"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/storage"
	"tgmessenger/internal/telegram"
)

// Run выполняет основную логику бота
func Run(ctx context.Context, client *telegram.Client) error {
	log.Println("🚀 Проверяем новые сообщения...")

	user, err := auth.AuthorizeClient(ctx, client.RawClient())
	if err != nil {
		return err
	}
	log.Printf("✅ Авторизован как %s", user.Username)

	unreadMessages, err := telegram.GetUnreadMessages(ctx, client)
	if err != nil {
		log.Printf("⚠️ Ошибка получения непрочитанных сообщений: %v", err)
	}

	directMessages, err := telegram.GetDirectMessages(ctx, client)
	if err != nil {
		log.Printf("⚠️ Ошибка поиска входящих сообщений: %v", err)
	}

	// Объединяем результаты двух методов
	messages := append(unreadMessages, directMessages...)

	if err := storage.SaveMessagesToJSON(messages); err != nil {
		log.Printf("⚠️ Ошибка сохранения JSON: %v", err)
	}

	if err := storage.SaveMessagesToMarkdown(messages); err != nil {
		log.Printf("⚠️ Ошибка сохранения Markdown: %v", err)
	}

	log.Println("✅ Обновление завершено.")
	return nil
}


