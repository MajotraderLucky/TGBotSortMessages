package auth

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

// AuthorizeClient выполняет авторизацию в Telegram
func AuthorizeClient(ctx context.Context, client *telegram.Client) (*tg.User, error) {
	fmt.Println("🔍 auth.AuthorizeClient: Старт") // Отладка

	api := tg.NewClient(client)
	authClient := client.Auth()

	// Проверяем статус авторизации
	status, err := authClient.Status(ctx)
	if err != nil {
		log.Printf("⚠️ Ошибка проверки статуса авторизации: %v", err)
		return nil, err
	}
	fmt.Printf("🔍 auth.AuthorizeClient: Статус авторизации: %v\n", status.Authorized) // Отладка

	// Если не авторизованы — входим
	if !status.Authorized {
		fmt.Println("🔍 auth.AuthorizeClient: Запуск авторизации...") // Отладка
		if err := Login(ctx, authClient); err != nil {
			log.Printf("⚠️ Ошибка авторизации: %v", err)
			return nil, err
		}
		fmt.Println("✅ auth.AuthorizeClient: Авторизация прошла") // Отладка
	}

	// Получаем информацию о текущем пользователе
	fmt.Println("🔍 auth.AuthorizeClient: Запрос информации о пользователе...") // Отладка
	users, err := api.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUserSelf{}})
	if err != nil {
		log.Printf("⚠️ Ошибка получения информации о себе: %v", err)
		return nil, err
	}

	user, ok := users[0].(*tg.User)
	if !ok {
		log.Println("⚠️ Ошибка приведения типа к *tg.User")
		return nil, fmt.Errorf("ошибка приведения типа к *tg.User")
	}

	fmt.Printf("✅ auth.AuthorizeClient: Авторизован как %s\n", user.Username) // Отладка
	return user, nil
}

