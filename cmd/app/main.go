package main

import (
	"log"

	"tgmessenger/internal/config"
	"tgmessenger/internal/telegram"
)

func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Запускаем бота через пакет telegram
	if err := telegram.Run(cfg); err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

