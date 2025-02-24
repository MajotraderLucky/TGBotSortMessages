package telegram

import (
	"github.com/gotd/td/session"
	"github.com/gotd/td/telegram"
	"tgmessenger/internal/config"
)

// NewClient создает новый клиент Telegram
func NewClient(cfg *config.Config) *telegram.Client {
	sessStorage := &session.FileStorage{Path: "session.json"}
	return telegram.NewClient(cfg.APIID, cfg.APIHash, telegram.Options{
		SessionStorage: sessStorage,
	})
}

