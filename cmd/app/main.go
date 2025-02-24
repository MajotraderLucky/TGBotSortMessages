package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gotd/td/telegram"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/session"
	"github.com/joho/godotenv"
)

// Загружаем конфигурацию из .env
func loadConfig() (int, string, error) {
	if err := godotenv.Load(); err != nil {
		return 0, "", err
	}

	apiID, err := strconv.Atoi(os.Getenv("API_ID"))
	if err != nil {
		return 0, "", err
	}

	apiHash := os.Getenv("API_HASH")
	if apiHash == "" {
		return 0, "", fmt.Errorf("API_HASH не указан")
	}

	return apiID, apiHash, nil
}

// Создаём клиент Telegram
func initTelegramClient(apiID int, apiHash string) *telegram.Client {
	sessStorage := &session.FileStorage{Path: "session.json"}
	return telegram.NewClient(apiID, apiHash, telegram.Options{
		SessionStorage: sessStorage,
	})
}

// Авторизация клиента
func authorizeClient(ctx context.Context, client *telegram.Client) (*tg.User, error) {
	api := tg.NewClient(client)
	authClient := client.Auth()

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
	api := tg.NewClient(client)

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
	apiID, apiHash, err := loadConfig()
	if err != nil {
		log.Fatal("Ошибка загрузки конфигурации:", err)
	}

	client := initTelegramClient(apiID, apiHash)
	ctx := context.Background()

	err = client.Run(ctx, func(ctx context.Context) error {
		log.Println("🚀 Бот запущен и подключен к Telegram")

		user, err := authorizeClient(ctx, client)
		if err != nil {
			return err
		}
		log.Printf("✅ Успешно авторизован как %s", user.Username)

		messages, err := getUnreadMessages(ctx, client)
		if err != nil {
			log.Printf("⚠️ Ошибка получения сообщений: %v", err)
			return nil
		}

		if err := saveMessagesToJSON(messages); err != nil {
			log.Printf("⚠️ Ошибка сохранения сообщений: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("Ошибка работы клиента: %v", err)
	}
}

