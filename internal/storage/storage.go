package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// SaveMessagesToJSON сохраняет сообщения в JSON
func SaveMessagesToJSON(messages []string) error {
	dir := "messages"
	filePath := filepath.Join(dir, "unread.json")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	log.Printf("📂 Сообщения сохранены в %s", filePath)
	return nil
}

// SaveStructuredMessagesToJSON сохраняет структурированные сообщения в JSON
func SaveStructuredMessagesToJSON(messages []map[string]interface{}) error {
	dir := "messages"
	filePath := filepath.Join(dir, "structured.json")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	log.Printf("📂 Структурированные сообщения сохранены в %s", filePath)
	return nil
}

// SaveSortedMessagesToJSON сохраняет сообщения, отсортированные по отправителям
func SaveSortedMessagesToJSON(messagesByUser map[string][]map[string]interface{}) error {
	dir := "messages"
	filePath := filepath.Join(dir, "sorted_by_user.json")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	data, err := json.MarshalIndent(messagesByUser, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	log.Printf("📂 Сортированные сообщения сохранены в %s", filePath)
	return nil
}

// SaveMessagesToMarkdown сохраняет 50 последних сообщений в .md
func SaveMessagesToMarkdown(messages []string) error {
	dir := "messages"
	filePath := filepath.Join(dir, "latest_messages.md")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Берем только 50 последних сообщений
	if len(messages) > 50 {
		messages = messages[len(messages)-50:]
	}

	// Формируем Markdown-формат с пустой строкой после заголовка
	mdContent := "# 📩 Последние 50 сообщений\n\n"

	if len(messages) > 0 {
		// Добавляем пустую строку перед списком
		for _, msg := range messages {
			mdContent += fmt.Sprintf("- %s\n", msg)
		}
		// Добавляем пустую строку после списка
		mdContent += "\n"
	} else {
		mdContent += "Нет новых сообщений.\n"
	}

	if err := os.WriteFile(filePath, []byte(mdContent), 0644); err != nil {
		return err
	}

	log.Printf("📂 Сообщения сохранены в %s", filePath)
	return nil
}

// SaveStructuredMessagesToMarkdown создает MD файл с более подробной информацией о сообщениях
func SaveStructuredMessagesToMarkdown(messages []map[string]interface{}) error {
	dir := "messages"
	filePath := filepath.Join(dir, "structured_messages.md")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Сортируем сообщения по дате (от новых к старым)
	sortFunc := func(i, j int) bool {
		// Получаем дату каждого сообщения
		dateI, okI := messages[i]["date"].(int32)
		dateJ, okJ := messages[j]["date"].(int32)

		// Если не можем получить даты, считаем сообщения равными
		if !okI || !okJ {
			return false
		}

		// Сортируем от новых к старым (в обратном порядке)
		return dateI > dateJ
	}

	messagesForMD := make([]map[string]interface{}, len(messages))
	copy(messagesForMD, messages)
	sort.Slice(messagesForMD, sortFunc)

	// Берем только 50 последних сообщений
	if len(messagesForMD) > 50 {
		messagesForMD = messagesForMD[:50]
	}

	// Формируем Markdown-формат с пустой строкой после заголовка
	mdContent := "# 📩 Последние 50 сообщений (отсортированные по времени)\n\n"

	if len(messagesForMD) > 0 {
		for _, msg := range messagesForMD {
			// Форматируем дату
			var dateStr string
			if date, ok := msg["date"].(int32); ok {
				dateStr = time.Unix(int64(date), 0).Format("2006-01-02 15:04:05")
			} else {
				dateStr = "неизвестно"
			}

			// Получаем текст сообщения и ID отправителя
			text, _ := msg["text"].(string)
			fromID, _ := msg["fromID"].(string)

			// Форматируем строку с пустыми строками до и после заголовка и текста
			mdContent += fmt.Sprintf("\n### От: %s (%s)\n\n", fromID, dateStr)
			mdContent += fmt.Sprintf("%s\n\n", text)

			// Добавляем разделитель
			mdContent += "---\n\n"
		}
	} else {
		mdContent += "Нет новых сообщений.\n"
	}

	if err := os.WriteFile(filePath, []byte(mdContent), 0644); err != nil {
		return err
	}

	log.Printf("📂 Структурированные сообщения сохранены в %s", filePath)
	return nil
}

// SaveUserGroupedMessagesToMarkdown сохраняет сообщения, сгруппированные по отправителям
func SaveUserGroupedMessagesToMarkdown(messagesByUser map[string][]map[string]interface{}) error {
	dir := "messages"
	filePath := filepath.Join(dir, "messages_by_user.md")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	// Формируем Markdown-формат с пустой строкой после заголовка
	mdContent := "# 📩 Сообщения по пользователям\n\n"

	// Получаем список пользователей и сортируем их по количеству сообщений
	var users []string
	for user := range messagesByUser {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return len(messagesByUser[users[i]]) > len(messagesByUser[users[j]])
	})

	if len(users) > 0 {
		for _, user := range users {
			messages := messagesByUser[user]

			// Заголовок с пользователем и количеством сообщений + пустая строка
			mdContent += fmt.Sprintf("\n## 👤 %s (%d сообщений)\n\n", user, len(messages))

			// Ограничиваем количество сообщений для каждого пользователя
			msgLimit := 10
			if len(messages) < msgLimit {
				msgLimit = len(messages)
			}

			// Добавляем сообщения пользователя с пустыми строками до и после заголовка
			for i := 0; i < msgLimit; i++ {
				msg := messages[i]

				// Форматируем дату
				var dateStr string
				if date, ok := msg["date"].(int32); ok {
					dateStr = time.Unix(int64(date), 0).Format("2006-01-02 15:04:05")
				} else {
					dateStr = "неизвестно"
				}

				// Получаем текст сообщения
				text, _ := msg["text"].(string)

				// Форматируем строку
				mdContent += fmt.Sprintf("\n### %s\n\n", dateStr)
				mdContent += fmt.Sprintf("%s\n\n", text)
			}

			// Если у пользователя больше сообщений, добавляем сноску
			if len(messages) > msgLimit {
				mdContent += fmt.Sprintf("\n... и еще %d сообщений\n\n", len(messages)-msgLimit)
			}

			// Добавляем разделитель между пользователями
			mdContent += "---\n\n"
		}
	} else {
		mdContent += "Нет сообщений для отображения.\n"
	}

	if err := os.WriteFile(filePath, []byte(mdContent), 0644); err != nil {
		return err
	}

	log.Printf("📂 Сообщения по пользователям сохранены в %s", filePath)
	return nil
}
