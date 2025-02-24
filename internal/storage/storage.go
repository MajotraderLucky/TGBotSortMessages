package storage

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
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

	// Формируем Markdown-формат
	mdContent := "# 📩 Последние 50 сообщений\n\n"
	for _, msg := range messages {
		mdContent += fmt.Sprintf("- %s\n", msg)
	}

	if err := os.WriteFile(filePath, []byte(mdContent), 0644); err != nil {
		return err
	}

	log.Printf("📂 Сообщения сохранены в %s", filePath)
	return nil
}

