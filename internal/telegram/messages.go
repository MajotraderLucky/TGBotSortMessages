package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/tg"
)

// GetDirectMessages получает входящие сообщения от всех пользователей, включая новых
func GetDirectMessages(ctx context.Context, client *Client) ([]string, error) {
	api := tg.NewClient(client.RawClient())

	resp, err := api.MessagesSearch(ctx, &tg.MessagesSearchRequest{
		Q:      "",
		Peer:   &tg.InputPeerEmpty{},
		Limit:  50,
		Filter: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска сообщений: %w", err)
	}

	var messages []string

	switch result := resp.(type) {
	case *tg.MessagesMessages:
		for _, msg := range result.Messages {
			if message, ok := msg.(*tg.Message); ok && message.Message != "" {
				fromID := "unknown"
				if message.FromID != nil {
					fromID = fmt.Sprintf("%v", message.FromID)
				}
				messages = append(messages, fmt.Sprintf("[%s]: %s", fromID, message.Message))
			}
		}
	case *tg.MessagesMessagesSlice:
		for _, msg := range result.Messages {
			if message, ok := msg.(*tg.Message); ok && message.Message != "" {
				fromID := "unknown"
				if message.FromID != nil {
					fromID = fmt.Sprintf("%v", message.FromID)
				}
				messages = append(messages, fmt.Sprintf("[%s]: %s", fromID, message.Message))
			}
		}
	default:
		log.Printf("⚠️ Неожиданный тип ответа в MessagesSearch: %T", resp)
	}

	return messages, nil
}

// GetUnreadMessages получает непрочитанные сообщения из чатов и групп
func GetUnreadMessages(ctx context.Context, client *Client) ([]string, error) {
	api := tg.NewClient(client.RawClient())

	resp, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{},
		Limit:      50,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения диалогов: %w", err)
	}

	var messages []string

	switch result := resp.(type) {
	case *tg.MessagesDialogs:
		for _, dialog := range result.Dialogs {
			if d, ok := dialog.(*tg.Dialog); ok && d.UnreadCount > 0 {
				inputPeer := GetInputPeer(d.Peer)
				if inputPeer == nil {
					continue
				}

				historyResp, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer:  inputPeer,
					Limit: d.UnreadCount,
				})
				if err != nil {
					log.Printf("⚠️ Ошибка получения истории сообщений: %v", err)
					continue
				}

				var messagesHistory []tg.MessageClass

				switch history := historyResp.(type) {
				case *tg.MessagesMessages:
					messagesHistory = history.Messages
				case *tg.MessagesMessagesSlice:
					messagesHistory = history.Messages
				default:
					log.Printf("⚠️ Неожиданный тип ответа в MessagesGetHistory: %T", historyResp)
					continue
				}

				for _, msg := range messagesHistory {
					if message, ok := msg.(*tg.Message); ok && message.Message != "" {
						fromID := "unknown"
						if message.FromID != nil {
							fromID = fmt.Sprintf("%v", message.FromID)
						}
						messages = append(messages, fmt.Sprintf("[%s]: %s", fromID, message.Message))
					}
				}
			}
		}
	default:
		log.Printf("⚠️ Неожиданный тип ответа в MessagesGetDialogs: %T", resp)
	}

	return messages, nil
}

