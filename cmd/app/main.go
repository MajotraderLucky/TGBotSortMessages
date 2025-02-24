package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/gotd/td/tg"
	"tgmessenger/internal/config"
	"tgmessenger/internal/telegram"
)

// Авторизация клиента
func authorizeClient(ctx context.Context, client *telegram.Client) (*tg.User, error) {
	api := tg.NewClient(client.RawClient())
	authClient := client.RawClient().Auth()

	status, err := authClient.Status(ctx)
	if err != nil {
		return nil, err
	}

	if !status.Authorized {
		return nil, fmt.Errorf("пользователь не авторизован")
	}

	users, err := api.UsersGetUsers(ctx, []tg.InputUserClass{&tg.InputUserSelf{}})
	if err != nil {
		return nil, err
	}

	user, ok := users[0].(*tg.User)
	if !ok {
		return nil, fmt.Errorf("ошибка приведения типа к *tg.User")
	}

	return user, nil
}

// Фильтр для чатов и пользователей (игнорируем каналы)
func getInputPeer(peer tg.PeerClass) tg.InputPeerClass {
	switch p := peer.(type) {
	case *tg.PeerUser:
		return &tg.InputPeerUser{UserID: p.UserID}
	case *tg.PeerChat:
		return &tg.InputPeerChat{ChatID: p.ChatID}
	default:
		return nil // Исключаем каналы
	}
}

// Получаем непрочитанные сообщения (без каналов)
func getUnreadMessages(ctx context.Context, client *telegram.Client) ([]string, error) {
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

		inputPeer := getInputPeer(dialog.Peer)
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

		// 📌 Обрабатываем все возможные типы ответов
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

		// Записываем только текстовые сообщения
		for _, msg := range messages {
			if message, ok := msg.(*tg.Message); ok && message.Message != "" {
				unreadMessages = append(unreadMessages, fmt.Sprintf("[%d]: %s", message.FromID, message.Message))
			}
		}
	}

	return unreadMessages, nil
}

// Сохраняем сообщения в JSON
func saveMessagesToJSON(messages []string) error {
	dir := "messages"
	filePath := filepath.Join(dir, "unread.json")

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.MarshalIndent(messages, "", "  ")
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	log.Printf("📂 Сообщения сохранены в %s", filePath)
	return nil
}

// Основная функция
func main() {
	// Загружаем конфигурацию
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	// Создаём Telegram клиента через обёртку
	client := telegram.NewClient(cfg)
	ctx := context.Background()

	err = client.RawClient().Run(ctx, func(ctx context.Context) error {
		log.Println("🚀 Бот запущен и подключен к Telegram")

		// Авторизация
		user, err := authorizeClient(ctx, client)
		if err != nil {
			return err
		}
		log.Printf("✅ Успешно авторизован как %s", user.Username)

		// Получение сообщений
		messages, err := getUnreadMessages(ctx, client)
		if err != nil {
			log.Printf("⚠️ Ошибка получения сообщений: %v", err)
			return nil
		}

		// Сохранение сообщений
		if err := saveMessagesToJSON(messages); err != nil {
			log.Printf("⚠️ Ошибка сохранения сообщений: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

