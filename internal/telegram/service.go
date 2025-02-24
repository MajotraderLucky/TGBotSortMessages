package telegram

import (
	"context"
	"log"

	"tgmessenger/internal/auth"
	"tgmessenger/internal/config"
)

// Service — структура для работы с ботом
type Service struct {
	client *Client
}

// NewService создаёт новый сервис Telegram
func NewService(cfg *config.Config) *Service {
	return &Service{
		client: NewClient(cfg),
	}
}

// Run запускает Telegram-клиент
func (s *Service) Run() error {
	ctx := context.Background()

	return s.client.RawClient().Run(ctx, func(ctx context.Context) error {
		log.Println("🚀 Бот запущен и подключен к Telegram")

		// Авторизация
		user, err := auth.AuthorizeClient(ctx, s.client.RawClient())
		if err != nil {
			return err
		}

		log.Printf("✅ Успешно авторизован как %s", user.Username)

		return nil
	})
}

