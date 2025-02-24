package telegram

import (
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/session"
	"tgmessenger/internal/config"
)

// Client обёртка вокруг gotd/td клиента
type Client struct {
	telegramClient *telegram.Client
}

// NewClient создаёт и возвращает Telegram клиента
func NewClient(cfg *config.Config) *Client {
	sessStorage := &session.FileStorage{Path: "session.json"}
	return &Client{
		telegramClient: telegram.NewClient(cfg.APIID, cfg.APIHash, telegram.Options{
			SessionStorage: sessStorage,
		}),
	}
}

// RawClient возвращает оригинальный *telegram.Client
func (c *Client) RawClient() *telegram.Client {
	return c.telegramClient
}

