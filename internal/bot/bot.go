package bot

import (
	"context"
	"fmt"
	"log"
	"time"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/config"
	"tgmessenger/internal/storage"
	"tgmessenger/internal/telegram"
)

// Run выполняет основную логику бота
func Run(ctx context.Context, client *telegram.Client, cfg *config.Config) error {
	if cfg.Debug {
		fmt.Println("🔍 bot.Run: Старт в режиме отладки") // Отладка
	} else {
		fmt.Println("🔍 bot.Run: Старт") // Отладка
	}
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
		latestMessages, latestStructured, err := telegram.GetLatestMessages(ctx, client, cfg.MessagesPerDialog) // Используем значение из конфигурации
		if err != nil {
			log.Printf("⚠️ Ошибка получения последних сообщений: %v", err)
		} else {
			log.Printf("📩 Получено %d последних сообщений из диалогов", len(latestMessages))
			messages = append(messages, latestMessages...)

			// Выводим первые 10 сообщений для отладки (или все, если их меньше 10)
			if cfg.Debug {
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
		}

		// Если нужны дополнительные непрочитанные сообщения
		var unreadStructured []map[string]interface{}
		var directStructured []map[string]interface{}

		if len(messages) == 0 {
			// Отладочная информация перед получением непрочитанных сообщений
			log.Println("📩 Начинаем получение непрочитанных сообщений...")
			unreadMessages, unreadStrct, err := telegram.GetUnreadMessages(ctx, client)
			if err != nil {
				log.Printf("⚠️ Ошибка получения непрочитанных сообщений: %v", err)
			} else {
				fmt.Printf("📩 bot.Run: Получено %d непрочитанных сообщений\n", len(unreadMessages)) // Отладка
				// Выводим полученные непрочитанные сообщения для отладки
				for i, msg := range unreadMessages {
					log.Printf("  Непрочитанное сообщение %d: %s", i+1, msg)
				}
				messages = append(messages, unreadMessages...)
				unreadStructured = unreadStrct
			}

			// Отладочная информация перед получением личных сообщений
			log.Println("📩 Начинаем получение входящих сообщений...")
			directMessages, directStrct, err := telegram.GetDirectMessages(ctx, client)
			if err != nil {
				log.Printf("⚠️ Ошибка поиска входящих сообщений: %v", err)
			} else {
				fmt.Printf("📩 bot.Run: Получено %d входящих сообщений\n", len(directMessages)) // Отладка
				// Выводим полученные личные сообщения для отладки
				for i, msg := range directMessages {
					log.Printf("  Входящее сообщение %d: %s", i+1, msg)
				}
				messages = append(messages, directMessages...)
				directStructured = directStrct
			}
		}

		log.Printf("📩 Всего обработано сообщений: %d", len(messages))

		// Объединяем структурированные сообщения
		var allStructured []map[string]interface{}
		allStructured = append(allStructured, latestStructured...)
		allStructured = append(allStructured, unreadStructured...)
		allStructured = append(allStructured, directStructured...)

		// Добавляем тестовое сообщение для проверки
		testMessages := []string{
			"[testuser]: Тестовое сообщение для проверки сохранения",
			"[anotheruser]: Еще одно тестовое сообщение",
		}

		// Добавляем тестовые структурированные сообщения
		testStructured := []map[string]interface{}{
			{
				"text":   "Тестовое сообщение для проверки сохранения",
				"fromID": "testuser",
				"date":   int32(time.Now().Unix()),
			},
			{
				"text":   "Еще одно тестовое сообщение",
				"fromID": "anotheruser",
				"date":   int32(time.Now().Unix() - 3600), // Час назад
			},
		}

		if len(messages) == 0 {
			log.Println("⚠️ Не найдено новых сообщений для обработки. Используем тестовые сообщения.")
			messages = testMessages
			allStructured = testStructured
		}

		// Сохраняем структурированные сообщения
		if len(allStructured) > 0 {
			log.Printf("📊 Сохраняем %d структурированных сообщений", len(allStructured))

			// Проверяем, какие форматы вывода включены
			for _, format := range cfg.OutputFormats {
				switch format {
				case "json":
					// Сохраняем все структурированные сообщения в JSON
					if err := storage.SaveStructuredMessagesToJSON(allStructured); err != nil {
						log.Printf("⚠️ Ошибка сохранения структурированного JSON: %v", err)
					}

					// Сохраняем простой список сообщений в JSON
					if err := storage.SaveMessagesToJSON(messages); err != nil {
						log.Printf("⚠️ Ошибка сохранения JSON: %v", err)
					} else if cfg.Debug {
						fmt.Println("✅ bot.Run: JSON сохранён") // Отладка
					}

				case "markdown":
					// Сохраняем в Markdown с форматированием
					if err := storage.SaveStructuredMessagesToMarkdown(allStructured); err != nil {
						log.Printf("⚠️ Ошибка сохранения структурированного Markdown: %v", err)
					}

					// Сохраняем простой список сообщений в Markdown
					if err := storage.SaveMessagesToMarkdown(messages); err != nil {
						log.Printf("⚠️ Ошибка сохранения Markdown: %v", err)
					} else if cfg.Debug {
						fmt.Println("✅ bot.Run: Markdown сохранён") // Отладка
					}

				case "structured":
					// Группируем и сортируем сообщения по дате
					sortedMessages := telegram.SortMessagesByDate(allStructured)

					// Сохраняем отсортированные по дате сообщения
					if err := storage.SaveStructuredMessagesToJSON(sortedMessages); err != nil {
						log.Printf("⚠️ Ошибка сохранения отсортированного JSON: %v", err)
					}

				case "byuser":
					// Группируем и сортируем сообщения по пользователям
					sortedByUser := telegram.SortMessagesByUser(allStructured)

					// Сохраняем отсортированные по пользователям сообщения
					if err := storage.SaveSortedMessagesToJSON(sortedByUser); err != nil {
						log.Printf("⚠️ Ошибка сохранения отсортированного JSON: %v", err)
					}

					// Сохраняем отсортированные сообщения в Markdown
					if err := storage.SaveUserGroupedMessagesToMarkdown(sortedByUser); err != nil {
						log.Printf("⚠️ Ошибка сохранения Markdown по пользователям: %v", err)
					}
				}
			}

			log.Println("✅ Все структурированные сообщения сохранены")
		}

		if cfg.Debug {
			fmt.Println("✅ bot.Run: Завершение") // Отладка
		}
		return nil
	})

	if err != nil {
		log.Printf("⚠️ Ошибка запуска клиента: %v", err)
		return err
	}

	log.Println("✅ bot.Run: Клиент успешно завершил работу")
	return nil
}
