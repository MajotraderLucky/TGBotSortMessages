package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"tgmessenger/internal/app"
	"tgmessenger/internal/auth"
)

func main() {
	fmt.Println("🔐 Программа авторизации в Telegram")

	client, ctx, cancel := app.InitBot()
	defer cancel()

	// Перехватываем Ctrl+C
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		fmt.Println("\n⛔ Прерывание пользователем. Выход...")
		cancel()
		os.Exit(1)
	}()

	// Получаем клиента Telegram
	rawClient := client.RawClient()

	// Запускаем клиента и выполняем авторизацию в его контексте
	log.Println("📡 Запускаем клиент Telegram...")
	err := rawClient.Run(ctx, func(ctx context.Context) error {
		// Выполняем вход, запрашивая номер телефона и код
		authClient := rawClient.Auth()
		log.Println("📱 Начинаем авторизацию...")
		if err := auth.Login(ctx, authClient); err != nil {
			return fmt.Errorf("ошибка авторизации: %w", err)
		}

		// Проверяем, что авторизация прошла успешно
		log.Println("🔍 Проверяем статус авторизации...")
		user, err := auth.AuthorizeClient(ctx, rawClient)
		if err != nil {
			return fmt.Errorf("ошибка проверки авторизации: %w", err)
		}

		log.Printf("✅ Авторизация успешна! Вы вошли как %s", user.Username)
		return nil
	})

	if err != nil {
		log.Fatalf("❌ Ошибка: %v", err)
	}

	fmt.Println("🔄 Теперь вы можете запустить основное приложение")
}
