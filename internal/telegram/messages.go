package telegram

import (
	"context"
	"fmt"
	"log"

	"github.com/gotd/td/tg"
)

// extractMessageText извлекает текст сообщения и ID отправителя
func extractMessageText(msg tg.MessageClass) (string, bool) {
	if message, ok := msg.(*tg.Message); ok && message.Message != "" {
		fromID := "unknown"
		if message.FromID != nil {
			fromID = fmt.Sprintf("%v", message.FromID)
		}
		return fmt.Sprintf("[%s]: %s", fromID, message.Message), true
	}
	return "", false
}

// processMessages обрабатывает список сообщений
func processMessages(messages []tg.MessageClass) []string {
	var result []string
	for _, msg := range messages {
		if text, ok := extractMessageText(msg); ok {
			result = append(result, text)
		}
	}
	return result
}

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
		messages = processMessages(result.Messages)
	case *tg.MessagesMessagesSlice:
		messages = processMessages(result.Messages)
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

				// Используем общую функцию для обработки сообщений
				messagesFromChat := processMessages(messagesHistory)
				messages = append(messages, messagesFromChat...)
			}
		}
	default:
		log.Printf("⚠️ Неожиданный тип ответа в MessagesGetDialogs: %T", resp)
	}

	return messages, nil
}
