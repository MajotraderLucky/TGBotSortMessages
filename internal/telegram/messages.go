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

// GetLatestMessages получает последние N сообщений из каждого диалога
func GetLatestMessages(ctx context.Context, client *Client, messagesPerDialog int) ([]string, error) {
	log.Printf("🔍 Запрашиваем последние %d сообщений из каждого диалога...", messagesPerDialog)

	api := tg.NewClient(client.RawClient())

	// Получаем список диалогов
	resp, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{},
		Limit:      100, // Увеличиваем лимит для получения большего количества диалогов
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения диалогов: %w", err)
	}

	log.Printf("✅ Получен ответ от MessagesGetDialogs: %T", resp)
	var allMessages []string
	var dialogCount int
	var processedDialogs int

	switch result := resp.(type) {
	case *tg.MessagesDialogs:
		dialogCount = len(result.Dialogs)
		log.Printf("📨 Найдено %d диалогов", dialogCount)

		for i, dialog := range result.Dialogs {
			d, ok := dialog.(*tg.Dialog)
			if !ok {
				continue
			}

			// Проверяем тип диалога
			switch d.Peer.(type) {
			case *tg.PeerUser:
				// Обрабатываем только личные диалоги
				log.Printf("📨 Обрабатываем личный диалог %d/%d", i+1, dialogCount)
			default:
				// Пропускаем группы и каналы
				log.Printf("📨 Пропускаем диалог %d/%d (не личный чат)", i+1, dialogCount)
				continue
			}

			inputPeer := GetInputPeer(d.Peer)
			if inputPeer == nil {
				log.Println("⚠️ Не удалось преобразовать Peer в InputPeer")
				continue
			}

			// Запрашиваем последние N сообщений
			historyResp, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
				Peer:  inputPeer,
				Limit: messagesPerDialog,
			})
			if err != nil {
				log.Printf("⚠️ Ошибка получения истории сообщений: %v", err)
				continue
			}

			processedDialogs++
			var messagesHistory []tg.MessageClass

			switch history := historyResp.(type) {
			case *tg.MessagesMessages:
				log.Printf("📨 Получено %d сообщений из диалога", len(history.Messages))
				messagesHistory = history.Messages
			case *tg.MessagesMessagesSlice:
				log.Printf("📨 Получено %d сообщений из диалога", len(history.Messages))
				messagesHistory = history.Messages
			case *tg.MessagesChannelMessages:
				log.Printf("📨 Получено %d сообщений из диалога", len(history.Messages))
				messagesHistory = history.Messages
			default:
				log.Printf("⚠️ Неожиданный тип ответа в MessagesGetHistory: %T", historyResp)
				continue
			}

			// Обрабатываем сообщения
			processedMessages := processMessages(messagesHistory)
			log.Printf("📨 Обработано %d сообщений из диалога", len(processedMessages))
			allMessages = append(allMessages, processedMessages...)
		}

	case *tg.MessagesDialogsSlice:
		dialogCount = len(result.Dialogs)
		log.Printf("📨 Найдено %d диалогов", dialogCount)

		for i, dialog := range result.Dialogs {
			d, ok := dialog.(*tg.Dialog)
			if !ok {
				continue
			}

			// Проверяем тип диалога
			switch d.Peer.(type) {
			case *tg.PeerUser:
				// Обрабатываем только личные диалоги
				log.Printf("📨 Обрабатываем личный диалог %d/%d", i+1, dialogCount)
			default:
				// Пропускаем группы и каналы
				log.Printf("📨 Пропускаем диалог %d/%d (не личный чат)", i+1, dialogCount)
				continue
			}

			inputPeer := GetInputPeer(d.Peer)
			if inputPeer == nil {
				log.Println("⚠️ Не удалось преобразовать Peer в InputPeer")
				continue
			}

			// Запрашиваем последние N сообщений
			historyResp, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
				Peer:  inputPeer,
				Limit: messagesPerDialog,
			})
			if err != nil {
				log.Printf("⚠️ Ошибка получения истории сообщений: %v", err)
				continue
			}

			processedDialogs++
			var messagesHistory []tg.MessageClass

			switch history := historyResp.(type) {
			case *tg.MessagesMessages:
				log.Printf("📨 Получено %d сообщений из диалога", len(history.Messages))
				messagesHistory = history.Messages
			case *tg.MessagesMessagesSlice:
				log.Printf("📨 Получено %d сообщений из диалога", len(history.Messages))
				messagesHistory = history.Messages
			case *tg.MessagesChannelMessages:
				log.Printf("📨 Получено %d сообщений из диалога", len(history.Messages))
				messagesHistory = history.Messages
			default:
				log.Printf("⚠️ Неожиданный тип ответа в MessagesGetHistory: %T", historyResp)
				continue
			}

			// Обрабатываем сообщения
			processedMessages := processMessages(messagesHistory)
			log.Printf("📨 Обработано %d сообщений из диалога", len(processedMessages))
			allMessages = append(allMessages, processedMessages...)
		}

	default:
		log.Printf("⚠️ Неожиданный тип ответа в MessagesGetDialogs: %T", resp)
	}

	log.Printf("📨 Всего получено %d сообщений из %d успешно обработанных диалогов (всего найдено: %d)",
		len(allMessages), processedDialogs, dialogCount)
	return allMessages, nil
}

// GetDirectMessages получает входящие сообщения от всех пользователей, включая новых
func GetDirectMessages(ctx context.Context, client *Client) ([]string, error) {
	log.Println("🔍 Запрашиваем входящие сообщения через MessagesSearch...")

	api := tg.NewClient(client.RawClient())

	// Попробуем получить диалоги сначала
	log.Println("🔍 Получаем список диалогов...")
	dialogsResp, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{},
		Limit:      50,
	})

	if err != nil {
		log.Printf("⚠️ Ошибка получения диалогов: %v", err)
	} else {
		log.Printf("✅ Получены диалоги: %T", dialogsResp)
	}

	// Теперь пробуем MessagesSearch
	resp, err := api.MessagesSearch(ctx, &tg.MessagesSearchRequest{
		Q:      "",
		Peer:   &tg.InputPeerEmpty{},
		Limit:  50,
		Filter: nil,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка поиска сообщений: %w", err)
	}

	log.Printf("✅ Получен ответ от MessagesSearch: %T", resp)
	var messages []string

	switch result := resp.(type) {
	case *tg.MessagesMessages:
		log.Printf("📨 MessagesMessages: найдено %d сообщений", len(result.Messages))
		messages = processMessages(result.Messages)
	case *tg.MessagesMessagesSlice:
		log.Printf("📨 MessagesMessagesSlice: найдено %d сообщений", len(result.Messages))
		messages = processMessages(result.Messages)
	case *tg.MessagesChannelMessages:
		log.Printf("📨 MessagesChannelMessages: найдено %d сообщений", len(result.Messages))
		messages = processMessages(result.Messages)
	case *tg.MessagesMessagesNotModified:
		log.Println("📨 MessagesMessagesNotModified: сообщения не изменились")
	default:
		log.Printf("⚠️ Неожиданный тип ответа в MessagesSearch: %T", resp)
	}

	return messages, nil
}

// GetUnreadMessages получает непрочитанные сообщения из чатов и групп
func GetUnreadMessages(ctx context.Context, client *Client) ([]string, error) {
	log.Println("🔍 Запрашиваем непрочитанные сообщения...")

	api := tg.NewClient(client.RawClient())

	resp, err := api.MessagesGetDialogs(ctx, &tg.MessagesGetDialogsRequest{
		OffsetPeer: &tg.InputPeerEmpty{},
		Limit:      50,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка получения диалогов: %w", err)
	}

	log.Printf("✅ Получен ответ от MessagesGetDialogs: %T", resp)
	var messages []string
	var dialogCount int

	switch result := resp.(type) {
	case *tg.MessagesDialogs:
		dialogCount = len(result.Dialogs)
		log.Printf("📨 MessagesDialogs: найдено %d диалогов", dialogCount)

		unreadDialogs := 0
		for _, dialog := range result.Dialogs {
			if d, ok := dialog.(*tg.Dialog); ok && d.UnreadCount > 0 {
				unreadDialogs++
				log.Printf("📨 Найден диалог с %d непрочитанными сообщениями", d.UnreadCount)

				inputPeer := GetInputPeer(d.Peer)
				if inputPeer == nil {
					log.Println("⚠️ Не удалось преобразовать Peer в InputPeer")
					continue
				}

				log.Printf("🔍 Запрашиваем историю сообщений для диалога...")
				historyResp, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer:  inputPeer,
					Limit: d.UnreadCount,
				})
				if err != nil {
					log.Printf("⚠️ Ошибка получения истории сообщений: %v", err)
					continue
				}

				log.Printf("✅ Получен ответ от MessagesGetHistory: %T", historyResp)
				var messagesHistory []tg.MessageClass

				switch history := historyResp.(type) {
				case *tg.MessagesMessages:
					log.Printf("📨 MessagesMessages: найдено %d сообщений в истории", len(history.Messages))
					messagesHistory = history.Messages
				case *tg.MessagesMessagesSlice:
					log.Printf("📨 MessagesMessagesSlice: найдено %d сообщений в истории", len(history.Messages))
					messagesHistory = history.Messages
				case *tg.MessagesChannelMessages:
					log.Printf("📨 MessagesChannelMessages: найдено %d сообщений в истории", len(history.Messages))
					messagesHistory = history.Messages
				default:
					log.Printf("⚠️ Неожиданный тип ответа в MessagesGetHistory: %T", historyResp)
					continue
				}

				// Используем общую функцию для обработки сообщений
				messagesFromChat := processMessages(messagesHistory)
				log.Printf("📨 Обработано %d сообщений из истории", len(messagesFromChat))
				messages = append(messages, messagesFromChat...)
			}
		}
		log.Printf("📨 Всего диалогов с непрочитанными сообщениями: %d", unreadDialogs)

	case *tg.MessagesDialogsSlice:
		dialogCount = len(result.Dialogs)
		log.Printf("📨 MessagesDialogsSlice: найдено %d диалогов", dialogCount)

		unreadDialogs := 0
		for _, dialog := range result.Dialogs {
			if d, ok := dialog.(*tg.Dialog); ok && d.UnreadCount > 0 {
				unreadDialogs++
				log.Printf("📨 Найден диалог с %d непрочитанными сообщениями", d.UnreadCount)

				inputPeer := GetInputPeer(d.Peer)
				if inputPeer == nil {
					log.Println("⚠️ Не удалось преобразовать Peer в InputPeer")
					continue
				}

				log.Printf("🔍 Запрашиваем историю сообщений для диалога...")
				historyResp, err := api.MessagesGetHistory(ctx, &tg.MessagesGetHistoryRequest{
					Peer:  inputPeer,
					Limit: d.UnreadCount,
				})
				if err != nil {
					log.Printf("⚠️ Ошибка получения истории сообщений: %v", err)
					continue
				}

				log.Printf("✅ Получен ответ от MessagesGetHistory: %T", historyResp)
				var messagesHistory []tg.MessageClass

				switch history := historyResp.(type) {
				case *tg.MessagesMessages:
					log.Printf("📨 MessagesMessages: найдено %d сообщений в истории", len(history.Messages))
					messagesHistory = history.Messages
				case *tg.MessagesMessagesSlice:
					log.Printf("📨 MessagesMessagesSlice: найдено %d сообщений в истории", len(history.Messages))
					messagesHistory = history.Messages
				case *tg.MessagesChannelMessages:
					log.Printf("📨 MessagesChannelMessages: найдено %d сообщений в истории", len(history.Messages))
					messagesHistory = history.Messages
				default:
					log.Printf("⚠️ Неожиданный тип ответа в MessagesGetHistory: %T", historyResp)
					continue
				}

				// Используем общую функцию для обработки сообщений
				messagesFromChat := processMessages(messagesHistory)
				log.Printf("📨 Обработано %d сообщений из истории", len(messagesFromChat))
				messages = append(messages, messagesFromChat...)
			}
		}
		log.Printf("📨 Всего диалогов с непрочитанными сообщениями: %d", unreadDialogs)

	default:
		log.Printf("⚠️ Неожиданный тип ответа в MessagesGetDialogs: %T", resp)
	}

	log.Printf("📨 Всего обработано %d непрочитанных сообщений из %d диалогов", len(messages), dialogCount)
	return messages, nil
}
