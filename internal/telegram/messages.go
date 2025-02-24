package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/tg"
)

// GetUnreadMessages получает непрочитанные сообщения (без каналов)
func GetUnreadMessages(ctx context.Context, client *Client) ([]string, error) {
	api := tg.NewClient(client.RawClient())

	resp, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{},
		Limit:      50,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения диалогов: %w", err)
	}

	var dialogs []tg.DialogClass
	switch d := resp.(type) {
	case *tg.MessagesDialogs:
		dialogs = d.Dialogs
	case *tg.MessagesDialogsSlice:
		dialogs = d.Dialogs
	default:
		return nil, fmt.Errorf("неожиданный тип ответа в MessagesGetDialogs: %T", resp)
	}

	var unreadMessages []string

	for _, dialogClass := range dialogs {
		dialog, ok := dialogClass.(*tg.Dialog)
		if !ok || dialog.UnreadCount == 0 {
			continue
		}

		inputPeer := GetInputPeer(dialog.Peer)
		if inputPeer == nil {
			continue // Пропускаем каналы
		}

		historyResp, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
			Peer:  inputPeer,
			Limit: dialog.UnreadCount,
		})
		if err != nil {
			log.Printf("⚠️ Ошибка получения истории сообщений: %v", err)
			continue
		}

		var messages []tg.MessageClass
		switch h := historyResp.(type) {
		case *tg.MessagesMessages:
			messages = h.Messages
		case *tg.MessagesMessagesSlice:
			messages = h.Messages
		case *tg.MessagesChannelMessages:
			messages = h.Messages
		default:
			log.Printf("⚠️ Неожиданный тип ответа в MessagesGetHistory: %T", historyResp)
			continue
		}

		for _, msg := range messages {
			if message, ok := msg.(*tg.Message); ok && message.Message != "" {
				fromID := "unknown"
				if message.FromID != nil {
					fromID = fmt.Sprintf("%v", message.FromID)
				}
				unreadMessages = append(unreadMessages, fmt.Sprintf("[%s]: %s", fromID, message.Message))
			}
		}
	}

	return unreadMessages, nil
}

