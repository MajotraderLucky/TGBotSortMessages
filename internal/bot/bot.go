package bot

import (
	"context"
	"fmt"
	"log"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/storage"
	"tgmessenger/internal/telegram"
)

// Run выполняет основную логику бота
func Run(ctx context.Context, client *telegram.Client) error {
	fmt.Println("🔍 bot.Run: Старт") // Отладка
	log.Println("🚀 Проверяем новые сообщения...")

	user, err := auth.AuthorizeClient(ctx, client.RawClient())
	if err != nil {
		log.Printf("⚠️ Ошибка авторизации: %v", err)
		return err
	}
	fmt.Println("✅ bot.Run: Авторизация прошла") // Отладка
	log.Printf("✅ Авторизован как %s", user.Username)

	unreadMessages, err := telegram.GetUnreadMessages(ctx, client)
	if err != nil {
		log.Printf("⚠️ Ошибка получения непрочитанных сообщений: %v", err)
	} else {
		fmt.Printf("📩 bot.Run: Получено %d непрочитанных сообщений\n", len(unreadMessages)) // Отладка
	}

	directMessages, err := telegram.GetDirectMessages(ctx, client)
	if err != nil {
		log.Printf("⚠️ Ошибка поиска входящих сообщений: %v", err)
	} else {
		fmt.Printf("📩 bot.Run: Получено %d входящих сообщений\n", len(directMessages)) // Отладка
	}

	// Объединяем результаты двух методов
	messages := append(unreadMessages, directMessages...)

	if err := storage.SaveMessagesToJSON(messages); err != nil {
		log.Printf("⚠️ Ошибка сохранения JSON: %v", err)
	} else {
		fmt.Println("✅ bot.Run: JSON сохранён") // Отладка
	}

	if err := storage.SaveMessagesToMarkdown(messages); err != nil {
		log.Printf("⚠️ Ошибка сохранения Markdown: %v", err)
	} else {
		fmt.Println("✅ bot.Run: Markdown сохранён") // Отладка
	}

	fmt.Println("✅ bot.Run: Завершение") // Отладка
	return nil
}

