package telegram

import (
	"context"
	"log"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/config"
)

// Run отвечает за запуск Telegram-клиента и авторизацию
func Run(cfg *config.Config) error {
	client := NewClient(cfg)
	ctx := context.Background()

	return client.RawClient().Run(ctx, func(ctx context.Context) error {
		log.Println("Бот запущен и подключен к Telegram")

		user, err := auth.AuthorizeClient(ctx, client.RawClient())
		if err != nil {
			return err
		}

		log.Printf("✅ Успешно авторизован как %s", user.Username)
		return nil
	})
}

