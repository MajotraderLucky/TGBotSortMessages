package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/session"
	"github.com/joho/godotenv"

	"tgmessenger/internal/auth"

)

func main() {
	// Загружаем переменные окружения
	if err := godotenv.Load(); err != nil {
		log.Fatal("Ошибка загрузки .env файла")
	}

	// Преобразуем API_ID в int
	apiID, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		log.Fatal("Ошибка конвертации API_ID в int:", err)
	}
	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		log.Fatal("API_HASH не указан в .env")
	}

	// Создаем хранилище сессии
	sessStorage := &session.FileStorage{Path: "session.json"}
	client := telegram.NewClient(apiID, apiHash, telegram.Options{
		SessionStorage: sessStorage,
	})

	ctx := context.Background()

	// Запускаем клиент
	err = client.Run(ctx, func(ctx context.Context) error {
		log.Println("Бот запущен и подключен к Telegram")
		api := tg.NewClient(client)

		// Проверяем авторизацию
		authClient := client.Auth()
		status, err := authClient.Status(ctx)
		if err != nil {
			return fmt.Errorf("Ошибка проверки статуса авторизации: %w", err)
		}
		if !status.Authorized {
			if err := auth.Login(ctx, authClient); err != nil { // Используем auth.Login
				return fmt.Errorf("Ошибка авторизации: %w", err)
			}
		}

		// Получаем информацию о текущем пользователе
		users, err := api.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUserSelf{}})
		if err != nil {
			return fmt.Errorf("Ошибка получения информации о себе: %w", err)
		}
		user, ok := users[0].(*tg.User)
		if !ok {
			return fmt.Errorf("Ошибка приведения типа к *tg.User")
		}
		log.Printf("✅ Успешно авторизован как %s", user.Username)

		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

