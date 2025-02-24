package storage

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

// SaveMessagesToJSON сохраняет сообщения в JSON файл
func SaveMessagesToJSON(messages []string) error {
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

