package auth

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
)

// AuthorizeClient выполняет авторизацию в Telegram
func AuthorizeClient(ctx context.Context, client *telegram.Client) (*tg.User, error) {
	api := tg.NewClient(client)
	authClient := client.Auth()

	// Проверяем статус авторизации
	status, err := authClient.Status(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка проверки статуса авторизации: %w", err)
	}

	// Если не авторизованы — входим
	if !status.Authorized {
		if err := Login(ctx, authClient); err != nil {
			return nil, fmt.Errorf("ошибка авторизации: %w", err)
		}
	}

	// Получаем информацию о текущем пользователе
	users, err := api.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUserSelf{}})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения информации о себе: %w", err)
	}
	user, ok := users[0].(*tg.User)
	if !ok {
		return nil, fmt.Errorf("ошибка приведения типа к *tg.User")
	}

	return user, nil
}

