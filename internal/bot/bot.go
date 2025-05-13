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

	// Получаем real клиент
	rawClient := client.RawClient()

	// Запускаем клиент
	err := rawClient.Run(ctx, func(ctx context.Context) error {
		log.Println("🚀 Запущен клиент в контексте...")

		// Авторизуемся
		user, err := auth.AuthorizeClient(ctx, rawClient)
		if err != nil {
			log.Printf("⚠️ Ошибка авторизации: %v", err)
			return err
		}
		fmt.Println("✅ bot.Run: Авторизация прошла") // Отладка
		log.Printf("✅ Авторизован как %s", user.Username)

		var messages []string

		// Попробуем получить последние сообщения из каждого диалога
		log.Println("📩 Получаем последние сообщения из диалогов...")
		latestMessages, err := telegram.GetLatestMessages(ctx, client, 3) // 3 последних сообщения из каждого диалога
		if err != nil {
			log.Printf("⚠️ Ошибка получения последних сообщений: %v", err)
		} else {
			log.Printf("📩 Получено %d последних сообщений из диалогов", len(latestMessages))
			messages = append(messages, latestMessages...)

			// Выводим первые 10 сообщений для отладки (или все, если их меньше 10)
			count := len(latestMessages)
			if count > 10 {
				count = 10
			}
			for i := 0; i < count; i++ {
				log.Printf("  Сообщение %d: %s", i+1, latestMessages[i])
			}
			if len(latestMessages) > 10 {
				log.Printf("  ... и ещё %d сообщений", len(latestMessages)-10)
			}
		}

		// Если нужны дополнительные непрочитанные сообщения
		if len(messages) == 0 {
			// Отладочная информация перед получением непрочитанных сообщений
			log.Println("📩 Начинаем получение непрочитанных сообщений...")
			unreadMessages, err := telegram.GetUnreadMessages(ctx, client)
			if err != nil {
				log.Printf("⚠️ Ошибка получения непрочитанных сообщений: %v", err)
			} else {
				fmt.Printf("📩 bot.Run: Получено %d непрочитанных сообщений\n", len(unreadMessages)) // Отладка
				// Выводим полученные непрочитанные сообщения для отладки
				for i, msg := range unreadMessages {
					log.Printf("  Непрочитанное сообщение %d: %s", i+1, msg)
				}
				messages = append(messages, unreadMessages...)
			}

			// Отладочная информация перед получением личных сообщений
			log.Println("📩 Начинаем получение входящих сообщений...")
			directMessages, err := telegram.GetDirectMessages(ctx, client)
			if err != nil {
				log.Printf("⚠️ Ошибка поиска входящих сообщений: %v", err)
			} else {
				fmt.Printf("📩 bot.Run: Получено %d входящих сообщений\n", len(directMessages)) // Отладка
				// Выводим полученные личные сообщения для отладки
				for i, msg := range directMessages {
					log.Printf("  Входящее сообщение %d: %s", i+1, msg)
				}
				messages = append(messages, directMessages...)
			}
		}

		log.Printf("📩 Всего обработано сообщений: %d", len(messages))

		// Добавляем тестовое сообщение для проверки
		testMessages := []string{
			"[testuser]: Тестовое сообщение для проверки сохранения",
			"[anotheruser]: Еще одно тестовое сообщение",
		}

		if len(messages) == 0 {
			log.Println("⚠️ Не найдено новых сообщений для обработки. Используем тестовые сообщения.")
			messages = testMessages
		}

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
	})

	if err != nil {
		log.Printf("⚠️ Ошибка запуска клиента: %v", err)
		return err
	}

	log.Println("✅ bot.Run: Клиент успешно завершил работу")
	return nil
}
